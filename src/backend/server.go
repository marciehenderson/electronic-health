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
	query string        // create, modify, view
	table string        // patient, location, record, etc.
	cols  []string      // patient_id, location_id, record_date, etc.
	keys  []string      // only required for view and modify, the primary keys to check against
	refs  []interface{} // only required for view and modify, the reference values to compare with primary keys
	data  []interface{} // only required for create and modify, can be any acceptable data type
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
			username := inputValidation(r.Form.Get("username"), "login")
			password := inputValidation(r.Form.Get("password"), "login")
			fmt.Println("Username: " + username)
			fmt.Println("Password: " + password)
			// Query database for user credentials
			if dbHandler(dbData{query: "view", table: "user", cols: []string{"password_hash"}, keys: []string{"user_hash"}, refs: []interface{}{username}}) == "[password_hash="+password+",]," {
				fmt.Println("Login successful")
				// Set cookies for user-accessible data
				value := make([]interface{}, 3)
				value[0] = dbHandler(dbData{query: "view", table: "record", cols: []string{"patient_id", "record_date", "location_id", "record_type", "notes", "created_at"}, keys: []string{""}, refs: []interface{}{""}})
				value[1] = dbHandler(dbData{query: "view", table: "client", cols: []string{"patient_id"}, keys: []string{"practitioner_id"}, refs: []interface{}{"placeholder"}})
				value[2] = dbHandler(dbData{query: "view", table: "user", cols: []string{"id"}, keys: []string{"user_hash"}, refs: []interface{}{username}})
				cookie := make([]http.Cookie, 3)
				cookie[0] = http.Cookie{Name: "record_data", Value: value[0].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
				cookie[1] = http.Cookie{Name: "client_data", Value: value[1].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
				cookie[2] = http.Cookie{Name: "user_data", Value: value[2].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
				http.SetCookie(w, &cookie[0])
				http.SetCookie(w, &cookie[1])
				http.SetCookie(w, &cookie[2])
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
			subHash := "/#actions+" + inputValidation(r.Form.Get("sub_hash"), "login")
			patientID := inputValidation(r.Form.Get("patient_id"), "login")
			locationID := inputValidation(r.Form.Get("location_id"), "login")
			recordDate := inputValidation(r.Form.Get("record_date"), "login")
			recordType := inputValidation(r.Form.Get("record_type"), "login")
			editValue := inputValidation(r.Form.Get("edit_value"), "login")
			notes := inputValidation(r.Form.Get("notes"), "login")
			fmt.Println("Sub Hash: " + subHash)
			fmt.Println("Patient ID: " + patientID)
			fmt.Println("Location ID: " + locationID)
			fmt.Println("Record Date: " + recordDate)
			fmt.Println("Record Type: " + recordType)
			fmt.Println("Edit Value: " + editValue)
			fmt.Println("Notes: " + notes)
			// Query database based on action sub-hash
			switch subHash {
			case "/#actions+create":
				dbHandler(dbData{query: "create", table: "record", cols: []string{"patient_id", "location_id", "record_type", "notes"}, data: []interface{}{patientID, locationID, recordType, notes}})
			case "/#actions+modify":
				dbHandler(dbData{query: "modify", table: "record", cols: []string{editValue}, keys: []string{"patient_id", "record_date"}, refs: []interface{}{patientID, recordDate}, data: []interface{}{"placeholder_repacement_value"}})
			case "/#actions+view":
				dbHandler(dbData{query: "view", table: "record", cols: []string{"patient_id", "record_date", "practitioner_id", "location_id", "notes", "code_cpt", "code_icd"}, keys: []string{"patient_id", "record_date"}, refs: []interface{}{patientID, recordDate}})
			}
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
	// Generate selection query based on data
	var selector string = ""
	var values string = " VALUES ("
	var inputs string = ""
	for i := 0; i < len(data.cols); i++ {
		// set selector, values, and inputs strings for each column in query
		selector += data.cols[i]
		if data.data != nil {
			// add quotes around string values
			if _, ok := data.data[i].(string); ok {
				inputs += "'" + data.data[i].(string) + "'"
			} else {
				inputs += data.data[i].(string)
			}
		}
		// add commas between values, but not at the end
		if i < len(data.cols)-1 {
			selector += ", "
			inputs += ", "
		}
	}
	inputs += ")"
	var comparator string = " WHERE "
	if data.keys != nil {
		for i := 0; i < len(data.keys); i++ {
			// set comparator string for each key in query
			comparator += data.keys[i] + " = '" + data.refs[i].(string) + "'"
			// add AND between keys, but not at the end
			if i < len(data.keys)-1 {
				comparator += " AND "
			}
		}
	}
	// Query selection successful
	fmt.Println("Query selection successful")
	// Query database
	switch data.query {
	case "create":
		// Generate insert query
		var insert string = "INSERT INTO " + data.table + " (" + selector + ") " + values + inputs
		// Create a new entry in the requested table
		_, err := db.Exec(insert)
		if err != nil {
			fmt.Println(err)
			fmt.Println(insert)
			return false
		}
	case "modify":
		// Generate update query
		var update string = "UPDATE " + data.table + " SET " + " (" + selector + ") " + values + inputs + comparator
		// Modify an existing entry in the requested table
		_, err := db.Exec(update)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
			fmt.Println(update)
			return false
		}
	case "view":
		// View an existing entry in the requested table
		output := ""
		colstr := make([]interface{}, len(data.cols))
		// Generate basic select query
		view := "SELECT " + selector + " FROM " + data.table
		// Add a WHERE clause if keys and refs are provided
		if data.keys[0] != "" {
			view += comparator
		}
		// Query the database for all relevant rows
		rows, err := db.Query(view)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
			return false
		}
		// Iterate through all rows and append to output
		// Output format: [col1=val1,col2=val2,...],[col1=val1,col2=val2,...],...
		for rows.Next() {
			rowstr := "["
			for i := 0; i < len(data.cols); i++ {
				colstr[i] = new(string)
			}
			// get value of each column
			err := rows.Scan(colstr...)
			if err != nil {
				fmt.Println(err)
				return false
			}
			fmt.Println("Scanned Row...")
			for i := 0; i < len(colstr); i++ {
				rowstr += data.cols[i] + "=" + *colstr[i].(*string) + ","
			}
			rowstr += "]"
			output += rowstr + ","
		}
		fmt.Println(output)
		// Return results
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
