package main

import (
	"net/http"

	"codeberg.org/textwire/textwire/v5"
	"codeberg.org/textwire/textwire/v5/config"
)

var tpl *textwire.Template

func main() {
	tpl, _ = textwire.NewTemplate(&config.Config{
		DebugMode:   true,
		FileWatcher: true,
	})

	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Response(w, "views/home", map[string]any{
		"title": "Home page",
	})

	err.PrintOnError()
}
