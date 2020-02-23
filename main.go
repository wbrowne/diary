package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var entries = []Entry{
	{ID: rand.Int63n(1000), Title: "I want to live forever", Text: "Palo santo gentrify next level deep v la croix pour-over. Vinyl gluten-free selfies trust fund hoodie next level raclette polaroid art party twee four dollar toast tbh microdosing green juice squid. Pickled pok pok crucifix selvage before they sold out heirloom kombucha retro truffaut salvia green juice coloring book stumptown portland fingerstache. Kale chips humblebrag gastropub fanny pack, taiyaki chia heirloom.\n\nArtisan thundercats deep v you probably haven't heard of them freegan twee tumeric snackwave dreamcatcher food truck cliche. Try-hard pickled drinking vinegar authentic affogato hella. Typewriter slow-carb shaman plaid ethical affogato. Bicycle rights snackwave schlitz twee semiotics man bun street art direct trade wolf 90's mlkshk selfies. Church-key taxidermy sriracha migas la croix kombucha. Chicharrones church-key dreamcatcher meditation paleo, echo park fashion axe gochujang raw denim chambray.", Created: time.Now(), Image: "https://images.unsplash.com/photo-1534774251706-35764d143b56"},
	{ID: rand.Int63n(1000), Title: "What a great day!", Text: "Palo santo gentrify next level deep v la croix pour-over. Vinyl gluten-free selfies trust fund hoodie next level raclette polaroid art party twee four dollar toast tbh microdosing green juice squid. Pickled pok pok crucifix selvage before they sold out heirloom kombucha retro truffaut salvia green juice coloring book stumptown portland fingerstache. Kale chips humblebrag gastropub fanny pack, taiyaki chia heirloom.\n\nArtisan thundercats deep v you probably haven't heard of them freegan twee tumeric snackwave dreamcatcher food truck cliche. Try-hard pickled drinking vinegar authentic affogato hella. Typewriter slow-carb shaman plaid ethical affogato. Bicycle rights snackwave schlitz twee semiotics man bun street art direct trade wolf 90's mlkshk selfies. Church-key taxidermy sriracha migas la croix kombucha. Chicharrones church-key dreamcatcher meditation paleo, echo park fashion axe gochujang raw denim chambray.", Created: time.Now(), Image: "https://images.unsplash.com/photo-1534774251706-35764d143b56"},
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "" // not required for 'postgres' admin user
	dbname   = "diary"
	enableDB = false
)

type Entry struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	Image   string    `json:"img"`
}

func main() {
	if enableDB {
		db, err := openDbConnection()
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected to Postgres")
	}
	setupServer()
}

func openDbConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db, err
}

func setupServer() {
	r := mux.NewRouter()
	r.HandleFunc("/entries", listEntries).Methods("GET", "OPTIONS")
	r.HandleFunc("/entries", createEntry).Methods("POST", "OPTIONS")
	r.HandleFunc("/entries/{id}", getEntry).Methods("GET", "OPTIONS")

	fmt.Println("Starting webserver @ 8080...")
	http.ListenAndServe(":8080", r)
}

func listEntries(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(entries)
}

func getEntry(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	requestedId, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range entries {
		if entry.ID == requestedId {
			json.NewEncoder(w).Encode(entry)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	var entry Entry
	json.NewDecoder(r.Body).Decode(&entry)
	entry.ID = rand.Int63n(1000)
	entry.Created = time.Now()
	entries = append(entries, entry)

	json.NewEncoder(w).Encode(&entry)
}
