package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", appHandler)
	fs := http.FileServer(http.Dir("../frontend/"))
	http.Handle("/..frontend/", http.StripPrefix("/..frontend/", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Handle GET requests - Only serve explicitly allowed files
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "../frontend/index.html")
		}
		if r.URL.Path == "/dist/bundle.js" {
			w.Header().Set("Content-Type", "application/javascript")
			http.ServeFile(w, r, "../frontend/dist/bundle.js")
		}
		if r.URL.Path == "/style.css" {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, "../frontend/style.css")
		}
	case "POST":
		// Handle POST requests - Not implemented
	default:
		// Handle all other requests - Not implemented
	}
}
