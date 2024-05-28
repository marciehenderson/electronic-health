package main

import (
	"fmt"
	"log"
	"net/http"
	"unicode"
)

func main() {
	http.HandleFunc("/", appHandler)
	fs := http.FileServer(http.Dir("../frontend/"))
	http.Handle("/..frontend/", http.StripPrefix("/..frontend/", fs))
	fmt.Println("Server running on: http://localhost:8080")
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
		if r.URL.Path == "/login" {
			// Parse form data
			r.ParseForm()
			var username = r.Form.Get("username")
			var password = r.Form.Get("password")
			// Validate user credentials
			username = inputValidation(username, "login")
			password = inputValidation(password, "login")
			fmt.Println("Username: " + username)
			fmt.Println("Password: " + password)
			// Check if user credentials are valid

			// reload account page
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "../frontend/index.html")
		}
	default:
		// Handle all other requests - Not implemented
	}
}

func inputValidation(input string, category string) string {
	// Validate input based on category
	var validatedInput string
	switch category {
	case "login":
		for _, c := range input {
			var character = string(c)
			// Sanitize character if it isn't alphanumeric
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				// convert to unicode
				character = "\\" + fmt.Sprintf("%U", c) + "\\"
			}
			validatedInput += character
		}
	default:
		fmt.Println("Invalid input category")
	}
	return validatedInput
}
