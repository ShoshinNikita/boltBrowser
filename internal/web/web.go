package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
	"github.com/ShoshinNikita/boltBrowser/internal/dbs"
)

// For embedding files
var (
	// "../../" for correct embedding of static files
	templates = packr.NewBox("../../templates")
	static    = packr.NewBox("../../static")
)

var routes = []struct {
	url     string
	method  string
	handler func(http.ResponseWriter, *http.Request)
}{
	{url: "/", method: "GET", handler: index},
	{url: "/wrapper", method: "GET", handler: wrapper},
	// databases
	{url: "/api/databases", method: "GET", handler: databasesList},
	{url: "/api/databases", method: "POST", handler: openDB},
	{url: "/api/databases", method: "DELETE", handler: closeDB},
	{url: "/api/databases/new", method: "POST", handler: createDB},
	// buckets
	{url: "/api/buckets", method: "POST", handler: addBucket},
	{url: "/api/buckets", method: "PUT", handler: editBucketName},
	{url: "/api/buckets", method: "DELETE", handler: deleteBucket},
	{url: "/api/buckets/current", method: "GET", handler: current},
	{url: "/api/buckets/root", method: "GET", handler: root},
	{url: "/api/buckets/back", method: "GET", handler: back},
	{url: "/api/buckets/next", method: "GET", handler: next},
	// records
	{url: "/api/records", method: "POST", handler: addRecord},
	{url: "/api/records", method: "PUT", handler: editRecord},
	{url: "/api/records", method: "DELETE", handler: deleteRecord},
	{url: "/api/records/prev", method: "GET", handler: prevRecords},
	{url: "/api/records/next", method: "GET", handler: nextRecords},
	// search
	{url: "/api/search", method: "GET", handler: search},
}

// Start website
func Start(port string, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	for _, r := range routes {
		router.Path(r.url).Methods(r.method).HandlerFunc(r.handler)
	}

	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static)))

	var handler http.Handler
	if config.Opts.Debug {
		handler = debugHandler(router)
	} else {
		handler = router
	}
	srv := http.Server{Addr: port, Handler: unescapingMiddleware(handler)}
	go srv.ListenAndServe()

	// Wait for signal
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
