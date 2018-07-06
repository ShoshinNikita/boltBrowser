package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ShoshinNikita/boltBrowser/internal/config"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("[ERR] %s\n", err.Error())
		fmt.Fprintf(w, "[ERR] %s\n", err.Error())
		return
	}
	data := struct{ WriteMode bool }{config.Opts.IsWriteMode}
	t.Execute(w, data)
}

func wrapper(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/wrapper.html")
	if err != nil {
		fmt.Printf("[ERR] %s\n", err.Error())
		fmt.Fprintf(w, "[ERR] %s\n", err.Error())
		return
	}

	data := map[string]interface{}{
		"URL": "http://localhost" + config.Opts.Port,
	}
	t.Execute(w, data)
}
