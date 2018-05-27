package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"dbs"
	"params"
)

var routes = []struct {
	url       string
	method    string
	handler   func(http.ResponseWriter, *http.Request)
	writeMode bool
}{
	{url: "/", method: "GET", handler: index},
	// databases
	{url: "/api/databases", method: "GET", handler: databasesList, writeMode: false},
	{url: "/api/databases", method: "POST", handler: openDB, writeMode: false},
	{url: "/api/databases", method: "DELETE", handler: closeDB, writeMode: false},
	// buckets
	{url: "/api/buckets", method: "POST", handler: addBucket, writeMode: true},
	{url: "/api/buckets", method: "DELETE", handler: deleteBucket, writeMode: true},
	{url: "/api/buckets/current", method: "GET", handler: current, writeMode: false},
	{url: "/api/buckets/root", method: "GET", handler: root, writeMode: false},
	{url: "/api/buckets/back", method: "GET", handler: back, writeMode: false},
	{url: "/api/buckets/next", method: "GET", handler: next, writeMode: false},
	// records
	{url: "/api/records", method: "POST", handler: addRecord, writeMode: true},
	{url: "/api/records", method: "PUT", handler: editRecord, writeMode: true},
	{url: "/api/records", method: "DELETE", handler: deleteRecord, writeMode: true},
	{url: "/api/records/prev", method: "GET", handler: prevRecords, writeMode: false},
	{url: "/api/records/next", method: "GET", handler: nextRecords, writeMode: false},
	// search
	{url: "/api/search", method: "GET", handler: search, writeMode: false},
}

// Start website
func Start(port string, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	for _, r := range routes {
		if !r.writeMode || (r.writeMode && params.IsWriteMode) {
			router.Path(r.url).Methods(r.method).HandlerFunc(r.handler)
		}
	}

	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	var handler http.Handler
	if params.Debug {
		handler = debugHandler(router)
	} else {
		handler = router
	}
	srv := http.Server{Addr: port, Handler: middleware(handler)}
	go srv.ListenAndServe()

	// Wait signal
	<-stopChan
	srv.Shutdown(context.Background())
	fmt.Println("[INFO] Website was stopped")
}

func debugHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		fmt.Printf("%s – %s\n", r.Method, r.URL.Path)
		if len(r.Form) > 0 {
			fmt.Print("Form:\n")
		}
		for key, values := range r.Form {
			fmt.Printf("* %s: %v\n", key, values)
		}

		fmt.Print("\n")

		h.ServeHTTP(w, r)
	})
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// Change all symbols
		for key, values := range r.Form {
			for i := range values {
				convertForBackend(&values[i])
			}
			r.Form[key] = values
		}

		h.ServeHTTP(w, r)
	})
}

func convertForBackend(origin *string) {
	s := *origin
	s = strings.Replace(s, "❮", "<", -1)
	s = strings.Replace(s, "❯", ">", -1)
	s = strings.Replace(s, "＂", "\"", -1)
	s = strings.Replace(s, "ߴ", "'", -1)
	*origin = s
}
