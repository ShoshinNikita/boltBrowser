package web

import (
	"os"
	//"fmt"
	"net/http"
	"log"
	
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/index.html")
}

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Methods("GET").HandlerFunc(index)
	
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))	

	log.Fatal(http.ListenAndServe(":500", handlers.LoggingHandler(os.Stdout, router)))
}