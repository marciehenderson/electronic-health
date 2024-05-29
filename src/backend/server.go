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
		// Handle POST requests
		if r.URL.Path == "/login" {
			// Parse login form data
			r.ParseForm()
			// Validate user credentials
			var username = inputValidation(r.Form.Get("username"), "login")
			var password = inputValidation(r.Form.Get("password"), "login")
			fmt.Println("Username: " + username)
			fmt.Println("Password: " + password)
			// reload account page
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "../frontend/index.html")
		}
		if r.URL.Path == "/action" {
			// Parse action form data
			r.ParseForm()
			var subHash = "/#actions+" + inputValidation(r.Form.Get("sub_hash"), "login")
			var patientID = inputValidation(r.Form.Get("patient_id"), "login")
			var locationID = inputValidation(r.Form.Get("location_id"), "login")
			var recordDate = inputValidation(r.Form.Get("record_date"), "login")
			var recordType = inputValidation(r.Form.Get("record_type"), "login")
			var editValue = inputValidation(r.Form.Get("edit_value"), "login")
			var notes = inputValidation(r.Form.Get("notes"), "login")
			fmt.Println("Sub Hash: " + subHash)
			fmt.Println("Patient ID: " + patientID)
			fmt.Println("Location ID: " + locationID)
			fmt.Println("Record Date: " + recordDate)
			fmt.Println("Record Type: " + recordType)
			fmt.Println("Edit Value: " + editValue)
			fmt.Println("Notes: " + notes)
			// reload action page
			w.Header().Set("Content-Type", "text/html")
			http.Redirect(w, r, subHash, http.StatusSeeOther)
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
