package web

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"db"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// Initialize – make(map[string]*db.DBApi)
func Initialize() {
	allDB = make(map[string]*db.BoltAPI)
}

// CloseDBs closes all databases
func CloseDBs() {
	for k := range allDB {
		allDB[k].Close()
		delete(allDB, k)
	}
	fmt.Println("[INFO] All databases were closed")
}

// Start runs website
func Start(port string, debug bool, stopChan chan struct{}) {
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/favicon.ico").Methods("GET").Handler(http.FileServer(http.Dir("./static/")))
	router.Path("/").Methods("GET").HandlerFunc(index)
	router.Path("/api/databases").Methods("POST").HandlerFunc(openDB)
	router.Path("/api/closeDB").Methods("POST").HandlerFunc(closeDB)
	router.Path("/api/databases").Methods("GET").HandlerFunc(databasesList)
	router.Path("/api/current").Methods("GET").HandlerFunc(current)
	router.Path("/api/root").Methods("GET").HandlerFunc(root)
	router.Path("/api/back").Methods("GET").HandlerFunc(back)
	router.Path("/api/next").Methods("GET").HandlerFunc(next)
	router.Path("/api/nextRecords").Methods("GET").HandlerFunc(nextRecords)
	router.Path("/api/prevRecords").Methods("GET").HandlerFunc(prevRecords)
	router.Path("/api/search").Methods("GET").HandlerFunc(search)
	router.Path("/api/searchRegex").Methods("GET").HandlerFunc(searchRegex)

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
