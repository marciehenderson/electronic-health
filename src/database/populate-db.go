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

// Declare global variables
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

	// Test name generation
	for i := 0; i < 10; i++ {
		name := randomName(3)
		fmt.Println(name[0], name[1])
		// save to log file
		err := os.WriteFile("log.txt", []byte(fmt.Sprintln(name[0], name[1])), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

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

func likelihood(dictionary []string, seed []string) []float32 {
	var weight []float32
	size := len(dictionary[0])
	for i := 0; i < len(dictionary); i++ {
		// check if all words in the dictionary are the same size
		if len(dictionary[i]) != size {
			return nil
		}
		var count int
		for j := 0; j < len(seed); j++ {
			for k := 0; k < len(seed[j])-size; k++ {
				// get a window of the seed from the current position
				window := seed[j][k : k+size]
				// check if the window matches the current word in the dictionary
				if window == dictionary[i] {
					count++
				}
			}
		}
		// calculate weight of the word based on its frequency in the seed
		weight = append(weight, float32(count)/float32(len(seed)))
	}
	return weight
}

// generates a random word with the requested number of syllables
// notably the word is not likely to be a real word, but should
// be useful for testing purposes.
func randomWord(syllables int) string {
	var word string
	dictionary := []string{"ba", "be", "bi", "bo", "bu", "da", "de", "di", "do", "du", "fa", "fe", "fi", "fo", "fu", "ga", "ge", "gi", "go", "gu", "ha", "he", "hi", "ho", "hu", "ja", "je", "ji", "jo", "ju", "ka", "ke", "ki", "ko", "ku", "la", "le", "li", "lo", "lu", "ma", "me", "mi", "mo", "mu", "na", "ne", "ni", "no", "nu", "pa", "pe", "pi", "po", "pu", "ra", "re", "ri", "ro", "ru", "sa", "se", "si", "so", "su", "ta", "te", "ti", "to", "tu", "va", "ve", "vi", "vo", "vu", "wa", "we", "wi", "wo", "wu", "ya", "ye", "yi", "yo", "yu", "za", "ze", "zi", "zo", "zu"}
	seed := []string{"alice", "bob", "charlie", "david", "eve", "frank", "grace", "heidi", "ivan", "judy", "kevin", "linda", "mallory", "nancy", "oscar", "peggy", "romeo", "sybil", "trudy", "victor", "walter", "xavier", "yvonne", "zelda"}
	weight := likelihood(dictionary, seed)
	// shuffle the dictionary and weight slices
	for i := range dictionary {
		j := rand.Intn(i + 1)
		dictionary[i], dictionary[j] = dictionary[j], dictionary[i]
		weight[i], weight[j] = weight[j], weight[i]
	}
	// make a random word with the requested number of syllables
	for i := 0; i < syllables; i++ {
		for j := 0; j < len(dictionary); j++ {
			if rand.Float32() < weight[j] {
				// concatenate a random syllable to the word
				word += dictionary[j]
				break
			}
		}
	}
	// print the random word to the console
	fmt.Println(word)
	return word
}

func randomName(max int) [2]string {
	var name [2]string
	// generate a random first name with a random number of syllables
	name[0] = randomWord(rand.Intn(max) + 1)
	// generate a random last name with a random number of syllables
	name[1] = randomWord(rand.Intn(max) + 1)
	return name
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
