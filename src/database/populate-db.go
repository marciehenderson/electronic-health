package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Get and set configuration values
	// Requires setting DBUSER and DBPASS environment variables in session
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "ehealth",
		AllowNativePasswords: true,
	}
	// Attempt to connect to the database
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	// Log any errors
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	// Log success
	fmt.Println("Connected to database")

	// Populate the database tables
	// User Table
	for i := 0; i < 100; i++ {
		// Generate a random user for the table
		_, err := db.Exec("INSERT INTO user (category, password_hash, user_hash, permissions) VALUES (?, ?, ?, ?)", randomCategory(), randomWord(10), randomWord(10), json.Number("0"))
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// Insert test user
	_, err = db.Exec("INSERT INTO user (category, password_hash, user_hash, permissions) VALUES (?, ?, ?, ?)", "admin", "test", "test", json.Number("1"))
	// Log any errors
	if err != nil {
		log.Fatal(err)
	}
	// Patient Table
	for i := 0; i < 100; i++ {
		// Generate a random patient for the table
		_, err := db.Exec("INSERT INTO patient (first_name, last_name, date_of_birth, street_address, contact_number, email) VALUES (?, ?, ?, ?, ?, ?)", randomWord(1), randomWord(1), randomDate(), randomWord(10), "1234567890", randomWord(5)+"@"+randomWord(5)+".com")
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// Practitioner Table
	for i := 0; i < 100; i++ {
		// Generate a random practitioner for the table
		_, err := db.Exec("INSERT INTO practitioner (first_name, last_name, date_of_birth, street_address, contact_number, email) VALUES (?, ?, ?, ?, ?, ?)", randomWord(1), randomWord(1), randomDate(), randomWord(10), "1234567890", randomWord(5)+"@"+randomWord(5)+".com")
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// Location Table
	for i := 0; i < 100; i++ {
		// Generate a random location for the table
		_, err := db.Exec("INSERT INTO location (location_name, street_address, contact_number, email) VALUES (?, ?, ?, ?)", randomWord(1), randomWord(10), "1234567890", randomWord(5)+"@"+randomWord(5)+".com")
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// Record Table
	for i := 0; i < 100; i++ {
		// Generate a random record for the table
		_, err := db.Exec("INSERT INTO record (patient_id, record_date, practitioner_id, location_id, notes, code_cpt, code_icd) VALUES (?, ?, ?, ?, ?, ?, ?)", rand.Intn(100)+1, randomDate(), rand.Intn(100)+1, rand.Intn(100)+1, randomWord(10), randomWord(5), randomWord(5))
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
	// Client Table
	for i := 0; i < 100; i++ {
		// Generate a random client for the table
		_, err := db.Exec("INSERT INTO client (patient_id, practitioner_id, client_status, visits) VALUES (?, ?, ?, ?)", rand.Intn(100)+1, rand.Intn(100)+1, "active", rand.Intn(1000))
		// Log any errors
		if err != nil {
			log.Fatal(err)
		}
	}
}

// generates a random word with the requested number of syllables
// notably the word is not likely to be a real word, but should
// be useful for testing purposes.
func randomWord(syllables int) string {
	var word string
	dictionary := []string{"ba", "be", "bi", "bo", "bu", "da", "de", "di", "do", "du", "fa", "fe", "fi", "fo", "fu", "ga", "ge", "gi", "go", "gu", "ha", "he", "hi", "ho", "hu", "ja", "je", "ji", "jo", "ju", "ka", "ke", "ki", "ko", "ku", "la", "le", "li", "lo", "lu", "ma", "me", "mi", "mo", "mu", "na", "ne", "ni", "no", "nu", "pa", "pe", "pi", "po", "pu", "ra", "re", "ri", "ro", "ru", "sa", "se", "si", "so", "su", "ta", "te", "ti", "to", "tu", "va", "ve", "vi", "vo", "vu", "wa", "we", "wi", "wo", "wu", "ya", "ye", "yi", "yo", "yu", "za", "ze", "zi", "zo", "zu"}
	// make a random word with the requested number of syllables
	for i := 0; i < syllables; i++ {
		// concatenate a random syllable to the word
		word += dictionary[rand.Intn(len(dictionary))]
	}
	return word
}

func randomCategory() string {
	categories := []string{"patient", "practitioner", "admin"}
	return categories[rand.Intn(len(categories))]
}

func randomDate() string {
	// random year between 1920 and 2019
	year := rand.Intn(100) + 1920
	// random month between 1 and 12
	month := rand.Intn(12) + 1
	// random day between 1 and 28
	day := rand.Intn(28) + 1
	// random hour between 0 and 23
	hour := rand.Intn(24)
	// random minute between 0 and 59
	minute := rand.Intn(60)
	// random second between 0 and 59
	second := rand.Intn(60)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}
