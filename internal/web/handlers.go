package web

import (
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/ShoshinNikita/log"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("").Parse(templates.String("index.html"))
	if err != nil {
		log.Errorf("%s\n", err.Error())
		fmt.Fprintf(w, "[ERR] %s\n", err.Error())
		return
	}
	t.Execute(w, nil)
}

// unescapingMiddleware use html.Unescape() for every element of r.Form
// For converting of "&lt;" into "<", "&gt;" into ">" and etc.
func unescapingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Errorf("Can't parse form: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for k := range r.Form {
			var values []string
			for _, s := range r.Form[k] {
				values = append(values, html.UnescapeString(s))
			}
			r.Form[k] = values
		}

		h.ServeHTTP(w, r)
	})
}
