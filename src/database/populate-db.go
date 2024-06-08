package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"database/populate-db/util"

	"github.com/go-sql-driver/mysql"
)

// Declare global variables
var db *sql.DB
var defaultSeed = []string{"apple", "bingo", "calendar", "dance", "evening", "flower", "goodbye", "hello", "ice", "jungle", "kite", "laptop", "mango", "notebook", "orange", "pencil", "queen", "rabbit", "sunset", "table", "umbrella", "violet", "water", "xylophone", "yellow", "zebra"}
var test = true // toggle test mode manually

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

	/**Testing***********************************************************************************************************************************************************************************************************/
	// Test name generation
	if test {
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		for i := 0; i < 100; i++ {
			name := randomName(3)
			// fmt.Println(name[0], name[1])
			// save to log file
			_, err := f.WriteString(fmt.Sprintf("%s %s\n", name[0], name[1]))
			if err != nil {
				log.Fatal(err)
			}
		}
		// Test util functions and types
		tree := util.GenerateWeightedTree(seedFromFile("name-seed.txt"))
		f, err = os.OpenFile("tree.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// Write the tree to a log file
		// Iterate over the dictionary and write each window and its weight to the file
		for i := 0; i < len(tree.Dictionary); i++ {
			_, err := f.WriteString(fmt.Sprintf("%s: %f\t[", tree.Dictionary[i].Window, tree.Weight[i]))
			if err != nil {
				log.Fatal(err)
			}
			// Iterate over the children of the current window and write each child and its weight to the file
			for j := 0; j < len(tree.Dictionary[i].Child); j++ {
				_, err := f.WriteString(fmt.Sprintf(" %s: %f,", *tree.Dictionary[i].Child[j].Window, tree.Dictionary[i].Child[j].Weight))
				if err != nil {
					log.Fatal(err)
				}
			}
			_, err = f.WriteString(" ]\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	/**Testing*End*******************************************************************************************************************************************************************************************************/

	// Populate the database tables
	// User Table
	if !test {
		for i := 0; i < 100; i++ {
			// Generate a random user for the table
			_, err := db.Exec("INSERT INTO user (category, password_hash, user_hash, permissions) VALUES (?, ?, ?, ?)", randomCategory(), randomWord(10, defaultSeed), randomWord(10, defaultSeed), json.Number("0"))
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
			_, err := db.Exec("INSERT INTO patient (first_name, last_name, date_of_birth, street_address, contact_number, email) VALUES (?, ?, ?, ?, ?, ?)", randomWord(1, defaultSeed), randomWord(1, defaultSeed), randomDate(), randomWord(10, defaultSeed), "1234567890", randomWord(5, defaultSeed)+"@"+randomWord(5, defaultSeed)+".com")
			// Log any errors
			if err != nil {
				log.Fatal(err)
			}
		}
		// Practitioner Table
		for i := 0; i < 100; i++ {
			// Generate a random practitioner for the table
			_, err := db.Exec("INSERT INTO practitioner (first_name, last_name, date_of_birth, street_address, contact_number, email) VALUES (?, ?, ?, ?, ?, ?)", randomWord(1, defaultSeed), randomWord(1, defaultSeed), randomDate(), randomWord(10, defaultSeed), "1234567890", randomWord(5, defaultSeed)+"@"+randomWord(5, defaultSeed)+".com")
			// Log any errors
			if err != nil {
				log.Fatal(err)
			}
		}
		// Location Table
		for i := 0; i < 100; i++ {
			// Generate a random location for the table
			_, err := db.Exec("INSERT INTO location (location_name, street_address, contact_number, email) VALUES (?, ?, ?, ?)", randomWord(1, defaultSeed), randomWord(10, defaultSeed), "1234567890", randomWord(5, defaultSeed)+"@"+randomWord(5, defaultSeed)+".com")
			// Log any errors
			if err != nil {
				log.Fatal(err)
			}
		}
		// Record Table
		for i := 0; i < 100; i++ {
			// Generate a random record for the table
			_, err := db.Exec("INSERT INTO record (patient_id, record_date, practitioner_id, location_id, notes, code_cpt, code_icd) VALUES (?, ?, ?, ?, ?, ?, ?)", rand.Intn(100)+1, randomDate(), rand.Intn(100)+1, rand.Intn(100)+1, randomWord(10, defaultSeed), randomWord(5, defaultSeed), randomWord(5, defaultSeed))
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
}

// Generates a weighted dictionary from a seed of strings with a given window size.
// The strength parameter determines how much the weight of a window is increased,
// and more specifically, the magnitude of difference between likely and unlikely
// windows. A strength of 0.0 will result in a uniform distribution of all windows.
func generateWeightedDictionary(seed []string, size int, strength float32) ([]string, []float32) {
	var dictionary []string
	var weight []float32
	for i := 0; i < len(seed); i++ {
		if len(seed[i]) < size {
			fmt.Printf("Error: seed string (%s) is too short for size of window (%d)\n", seed[i], size)
			break
		}
		for j := 0; j < len(seed[i])-size; j++ {
			// get a window of the seed from the current position
			window := seed[i][j : j+size]
			// check if the window is already in the dictionary
			found := false
			// only check the dictionary if it is not empty
			if len(dictionary) > 0 {
				for k := 0; k < len(dictionary); k++ {
					if window == dictionary[k] {
						// increase the weight of the window based on the strength
						weight[k] += strength
						found = true
						break
					}
				}
			}
			// add the window to the dictionary if it is not already there
			// and initialize its associated weight to 1.0
			if !found {
				dictionary = append(dictionary, window)
				weight = append(weight, 1.0)
			}
		}
	}
	return dictionary, weight
}

// Selects a window from a weighted dictionary based on the weight of each window.
func selectWindow(dictionary []string, weight []float32) string {
	// calculate the total weight of the dictionary
	var total float32
	for i := 0; i < len(weight); i++ {
		total += weight[i]
	}
	// generate a random number between 0 and the total weight
	rand := rand.Float32() * total
	// select a window based on its weight
	for i := 0; i < len(weight); i++ {
		if rand < weight[i] {
			return dictionary[i]
		}
		rand -= weight[i]
	}
	return ""
}

// Reads a file and returns its contents as a string array
func seedFromFile(path string) []string {
	// open the file at the given path
	f, err := os.Open(path)
	// log any errors
	if err != nil {
		log.Fatal(err)
	}
	// close the file when the function returns
	defer f.Close()
	// read the file into a byte slice with a maximum size of 3.84 MB
	data := make([]byte, 3840000)
	count, err := f.Read(data)
	// log any errors
	if err != nil {
		log.Fatal(err)
	}
	// convert data to []string with EOL as delimiter
	raw := string(data[:count])
	var seed []string
	for i := 0; i < len(raw); i++ {
		// check for the end of the line
		if raw[i] == '\n' {
			// append raw slice to seed
			seed = append(seed, raw[:i])
			// remove slice from raw
			raw = raw[i+1:]
			// reset iterator
			i = 0
		}
	}
	// return the byte slice as a string
	return seed
}

// Generates a random word with the requested number of syllables.
// Notably the word is not likely to be a real word, but should
// be useful for testing purposes.
func randomWord(syllables int, seed []string) string {
	var word string
	dictionary, weight := generateWeightedDictionary(seed, 3, 0.5)
	// make a random word with the requested number of syllables
	for i := 0; i < syllables; i++ {
		word += selectWindow(dictionary, weight)
	}
	// print the random word to the console
	//fmt.Println(word)
	return word
}

func randomName(max int) [2]string {
	//seed := []string{"alice", "bob", "charlie", "david", "eve", "frank", "grace", "heidi", "ivan", "judy", "kevin", "linda", "mallory", "nancy", "oscar", "peggy", "romeo", "sybil", "trudy", "victor", "walter", "xavier", "yvonne", "zelda"}
	seed := seedFromFile("name-seed.txt")
	var name [2]string
	// generate a random first name with a random number of syllables
	name[0] = randomWord(rand.Intn(max)+1, seed)
	// generate a random last name with a random number of syllables
	name[1] = randomWord(rand.Intn(max)+1, seed)
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
