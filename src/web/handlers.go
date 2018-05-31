package web

import (
	"fmt"
	"html/template"
	"net/http"

	"params"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("[ERR] %s\n", err.Error())
		fmt.Fprintf(w, "[ERR] %s\n", err.Error())
		return
	}
	data := struct{ WriteMode bool }{params.IsWriteMode}
	t.Execute(w, data)
}
