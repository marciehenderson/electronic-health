package main

// Import necessary packages
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/go-sql-driver/mysql"
)

// Declare global variables
var dbConfig = mysql.Config{
	User:                 os.Getenv("DBUSER"),
	Passwd:               os.Getenv("DBPASS"),
	Net:                  "tcp",
	Addr:                 "localhost:3306",
	DBName:               "ehealth",
	AllowNativePasswords: true,
}

// Declare custom types
type dbData struct {
	query  string        // create, modify, view
	table  string        // patient, location, record, etc.
	column string        // patient_id, location_id, record_date, etc.
	colRef string        // required for view
	rowRef []interface{} // only required for view and modify, can be multiple primary keys
	data   interface{}   // only required for create and modify, can be any acceptable data type
}

// Main function of the server
func main() {
	// Configure and start the server
	http.HandleFunc("/", appHandler)
	fs := http.FileServer(http.Dir("../frontend/"))
	http.Handle("/..frontend/", http.StripPrefix("/..frontend/", fs))
	fmt.Println("Server running on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle all requests to the server
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
			// Query database for user credentials
			if dbHandler(dbData{query: "view", table: "user", column: "password_hash", colRef: "user_hash", rowRef: []interface{}{username}}) == password {
				fmt.Println("Login successful")
			} else {
				fmt.Println("Login failed")
			}
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

// Handle all database queries
func dbHandler(data dbData) interface{} {
	// Connect to database
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	// Log any errors
	if err != nil {
		log.Fatal(err)
		return false
	}
	pingErr := db.Ping()
	if db.Ping() != nil {
		log.Fatal(pingErr)
		return false
	}
	// Connection successful
	fmt.Println("Connected to database")
	// Query database
	switch data.query {
	case "create":
		// Create a new entry in the requested table
		_, err := db.Exec("INSERT INTO "+data.table+" ("+data.column+") VALUES (?)", data.data)
		if err != nil {
			log.Fatal(err)
			return false
		}
	case "modify":
		// Modify an existing entry in the requested table
		_, err := db.Exec("UPDATE "+data.table+" SET "+data.column+" = ? WHERE "+data.column+" = ?", data.data, data.rowRef)
		if err != nil {
			log.Fatal(err)
			return false
		}
	case "view":
		// View an existing entry in the requested table
		var output string
		err := db.QueryRow("SELECT "+data.column+" FROM "+data.table+" WHERE "+data.colRef+" = ?", data.rowRef[0].(string)).Scan(&output)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
			return false
		}
		// Return results
		fmt.Print("Output: ")
		return output
	default:
		fmt.Println("Invalid query type")
		return false
	}
	// Return results
	return true
}

// Validate input string based on category
func inputValidation(input string, category string) string {
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
