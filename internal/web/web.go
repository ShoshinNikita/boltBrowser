package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"

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
func Start(port int, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	for _, r := range routes {
		router.Path(r.url).Methods(r.method).HandlerFunc(r.handler)
	}

	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static)))

	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: unescapingMiddleware(router)}
	go srv.ListenAndServe()

	// Wait for signal
	<-stopChan
	srv.Shutdown(context.Background())
}
