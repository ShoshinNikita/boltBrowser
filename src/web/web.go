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
	// recrods
	{url: "/api/records", method: "POST", handler: addRecord, writeMode: true},
	{url: "/api/records", method: "PUT", handler: editRecord, writeMode: true},
	{url: "/api/records", method: "DELETE", handler: deleteRecord, writeMode: true},
	{url: "/api/records/prev", method: "GET", handler: prevRecords, writeMode: false},
	{url: "/api/records/next", method: "GET", handler: nextRecords, writeMode: false},
	// search
	{url: "/api/search", method: "GET", handler: search, writeMode: false},
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
	}{writeMode}
	t.Execute(w, data)
}

// TODO
var writeMode = true

// Start runs website
func Start(port string, debug bool, stopChan chan struct{}) {
	dbs.Init()

	router := mux.NewRouter().StrictSlash(false)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	for _, r := range routes {
		if !r.writeMode || (r.writeMode && writeMode) {
			router.Path(r.url).Methods(r.method).HandlerFunc(r.handler)
		}
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
