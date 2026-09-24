package server

import (
	"fmt"
	"net/http"
	"os"

	"codeberg.org/textwire/textwire/v5"
	"codeberg.org/textwire/textwire/v5/config"
)

var tpl *textwire.Template

func Serve(port string) {
	setupTextwire()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", reposHandler)
	http.HandleFunc("/tags/{name}", tagsHandler)
	http.ListenAndServe(port, nil)
}

func setupTextwire() {
	isDev := os.Getenv("APP_ENV") == "development"

	twConf := &config.Config{
		DebugMode:     isDev,
		FileWatcher:   isDev,
		ErrorPagePath: "error-page",
		GlobalData: map[string]any{
			"isDev": isDev,
		},
	}

	t, err := textwire.NewTemplate(twConf)
	if err != nil {
		fmt.Println(err)
	}
	tpl = t
}
