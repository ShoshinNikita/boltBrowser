package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"dbs"
)

var routes = []struct {
	url     string
	method  string
	handler func(http.ResponseWriter, *http.Request)
}{
	// TODO next -> buckets/next
	// prevRecords -> records/prev etc.
	{url: "/", method: "GET", handler: index},
	{url: "/api/databases", method: "POST", handler: openDB},
	{url: "/api/closeDB", method: "POST", handler: closeDB},
	{url: "/api/databases", method: "GET", handler: databasesList},
	{url: "/api/current", method: "GET", handler: current},
	{url: "/api/root", method: "GET", handler: root},
	{url: "/api/back", method: "GET", handler: back},
	{url: "/api/next", method: "GET", handler: next},
	{url: "/api/nextRecords", method: "GET", handler: nextRecords},
	{url: "/api/prevRecords", method: "GET", handler: prevRecords},
	{url: "/api/search", method: "GET", handler: search},

	{url: "/api/buckets", method: "POST", handler: addBucket},
	{url: "/api/buckets", method: "DELETE", handler: deleteBucket},
	{url: "/api/records", method: "POST", handler: addRecord},
	{url: "/api/records", method: "PUT", handler: editRecord},
	{url: "/api/records", method: "DELETE", handler: deleteRecord},
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("[ERR] %s\n", err.Error())
		fmt.Fprintf(w, "[ERR] %s\n", err.Error())
		return
	}
	// TODO global config
	data := struct {
		WriteMode bool
	}{true}
	t.Execute(w, data)
}

// Start runs website
func Start(port string, debug bool, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	for _, r := range routes {
		router.Path(r.url).Methods(r.method).HandlerFunc(r.handler)
	}

	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	var handler http.Handler
	if debug {
		handler = handlers.LoggingHandler(os.Stdout, router)
	} else {
		handler = router
	}
	srv := http.Server{Addr: port, Handler: handler}
	go srv.ListenAndServe()

	// Wait signal
	<-stopChan
	srv.Shutdown(context.Background())
	fmt.Println("[INFO] Website was stopped")
}
