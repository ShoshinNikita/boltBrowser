package web

import (
	"os"
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"db"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// Initialize – make(map[string]*db.DBApi)
func Initialize() {
	allDB = make(map[string]*db.DBApi)
}

// Start runs website
func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Methods("GET").HandlerFunc(index)
	router.Path("/api/openDB").Methods("POST").HandlerFunc(openDB)
	router.Path("/api/closeDB").Methods("POST").HandlerFunc(closeDB)
	router.Path("/api/databases").Methods("GET").HandlerFunc(databasesList)
	router.Path("/api/current").Methods("GET").HandlerFunc(current)
	router.Path("/api/cmd").Methods("GET").HandlerFunc(cmd)
	router.Path("/api/back").Methods("GET").HandlerFunc(back)
	router.Path("/api/next").Methods("GET").HandlerFunc(next)

	// For static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))	

	log.Fatal(http.ListenAndServe(":500", handlers.LoggingHandler(os.Stdout, router)))
}