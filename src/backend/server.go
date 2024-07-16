package main

// Import necessary packages
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
			username := inputValidation(r.Form.Get("username"), "basic")
			password := inputValidation(r.Form.Get("password"), "basic")
			// fmt.Println("Username: " + username)
			// fmt.Println("Password: " + password)
			// Query database for user credentials
			if dbHandler(dbData{query: "view", table: "user", cols: []string{"password_hash"}, keys: []string{"user_hash"}, refs: []interface{}{username}}) == "[password_hash="+password+",]," {
				fmt.Println("Login successful")
				// Set cookies for user-accessible data
				value := make([]interface{}, 4)
				value[0] = dbHandler(dbData{query: "view", table: "record", cols: []string{"patient_id", "record_date", "location_id", "record_type", "notes", "created_at"}, keys: nil, refs: nil})
				value[1] = dbHandler(dbData{query: "view", table: "user", cols: []string{"id"}, keys: []string{"user_hash"}, refs: []interface{}{username}})
				// get user id from user data
				userID := strings.Split(strings.Split(value[1].(string), "id=")[1], ",")[0]
				// fmt.Println("User ID: " + userID)
				value[2] = dbHandler(dbData{query: "view", table: "client", cols: []string{"patient_id"}, keys: []string{"practitioner_id"}, refs: []interface{}{userID}})
				// get array of patient ids from client data
				firstCut := strings.Split(value[2].(string), "[patient_id=")
				patientIDs := make([]string, len(firstCut))
				for i := 0; i < len(firstCut); i++ {
					patientIDs[i] = strings.Split(firstCut[i], ",")[0]
				}
				// fmt.Println(patientIDs)
				value[3] = dbHandler(dbData{query: "view", table: "patient", cols: []string{"id", "first_name", "last_name", "date_of_birth", "street_address", "contact_number", "email", "created_at", "updated_at"}, keys: []string{"id"}, refs: []interface{}{patientIDs}})
				cookie := make([]http.Cookie, 4)
				_, ok := value[0].(string)
				if ok {
					cookie[0] = http.Cookie{Name: "record_data", Value: value[0].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
					http.SetCookie(w, &cookie[0])
				}
				_, ok = value[1].(string)
				if ok {
					cookie[1] = http.Cookie{Name: "user_data", Value: value[1].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
					http.SetCookie(w, &cookie[1])
				}
				_, ok = value[2].(string)
				if ok {
					cookie[2] = http.Cookie{Name: "client_data", Value: value[2].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
					http.SetCookie(w, &cookie[2])
				}
				_, ok = value[3].(string)
				if ok {
					cookie[3] = http.Cookie{Name: "patient_data", Value: value[3].(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
					http.SetCookie(w, &cookie[3])
				}
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
			subHash := "/#actions+" + inputValidation(r.Form.Get("sub_hash"), "basic")
			patientID := inputValidation(r.Form.Get("patient_id"), "basic")
			locationID := inputValidation(r.Form.Get("location_id"), "basic")
			recordDate := inputValidation(r.Form.Get("record_date"), "datetime")
			recordType := inputValidation(r.Form.Get("record_type"), "basic")
			notes := inputValidation(r.Form.Get("notes"), "basic")
			// Query database based on action sub-hash
			switch subHash {
			case "/#actions+create":
				// create new record
				dbHandler(dbData{query: "create", table: "record", cols: []string{"patient_id", "location_id", "record_type", "notes"}, data: []interface{}{patientID, locationID, recordType, notes}})
			case "/#actions+modify":
				// modify existing record
				// dbHandler(dbData{query: "modify", table: "record", cols: []string{editValue}, keys: []string{"patient_id", "record_date"}, refs: []interface{}{patientID, recordDate}, data: []interface{}{"placeholder_repacement_value"}})
			case "/#actions+view":
				// view existing record
				value := dbHandler(dbData{query: "view", table: "record", cols: []string{"patient_id", "record_date", "location_id", "record_type", "notes"}, keys: []string{"patient_id", "record_date"}, refs: []interface{}{patientID, recordDate}})
				// set cookie for currently requested record
				cookie := http.Cookie{Name: "record_view", Value: value.(string), Path: "/", SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: false}
				http.SetCookie(w, &cookie)
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
			// use IN for keys with mutiple refs
			if _, ok := data.refs[i].([]string); ok {
				comparator += data.keys[i] + " IN ("
				for j := 0; j < len(data.refs[i].([]string)); j++ {
					comparator += "'" + data.refs[i].([]string)[j] + "'"
					// add commas between values, but not at the end
					if j < len(data.refs[i].([]string))-1 {
						comparator += ", "
					}
				}
				comparator += ")"
			} else {
				// set comparator string for each key in query
				comparator += data.keys[i] + " = '" + data.refs[i].(string) + "'"
			}
			// add AND between keys, but not at the end
			if i < len(data.keys)-1 {
				comparator += " AND "
			}
		}
	}
	// Query selection successful
	// fmt.Println("Query selection successful")
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
		if data.keys != nil {
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
				// skip to next row if error
				continue
			}
			// fmt.Println("Scanned Row...")
			for i := 0; i < len(colstr); i++ {
				rowstr += data.cols[i] + "=" + *colstr[i].(*string) + ","
			}
			rowstr += "]"
			output += rowstr + ","
		}
		// fmt.Println(output)
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
	var validated string
	switch category {
	case "basic":
		// basic method of validating input.
		// simply sanitizes string to replace
		// non-alphanumeric with unicode.
		validated = sanitizeString(input)
	case "datetime":
		// input validation for datetime
		// inputs. ensures values are compliant
		// with SQL datatype structure.
		validated = sanitizeString(input)
		// allow hyphens in fetch requests
		// for compliance with DATE and
		// TIMESTAMP datatypes.
		validated = strings.ReplaceAll(validated, "\\U+002D\\", "-")
		// allow colons in fetch requests
		validated = strings.ReplaceAll(validated, "\\U+003A\\", ":")
		// allow spaces in fetch requests
		validated = strings.ReplaceAll(validated, "\\U+0020\\", " ")
	default:
		fmt.Println("Invalid input category")
	}
	return validated
}

func sanitizeString(input string) string {
	var output string
	for _, c := range input {
		var character = string(c)
		// sanitize non-alphanumeric characters
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			// convert character to unicode
			character = "\\" + fmt.Sprintf("%U", c) + "\\"
		}
		output += character
	}
	return output
}
