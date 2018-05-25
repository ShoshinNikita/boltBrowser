package web

import (
	"context"
	"fmt"
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
}

// Start runs website
func Start(port string, debug bool, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	router.Path("/").Methods("GET").Handler(http.FileServer(http.Dir("templates/")))
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
