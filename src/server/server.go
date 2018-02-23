package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Game This struct holds all the game data
type Game struct {
	ID    int64
	User  string
	Word  string
	Guess string
}

// GameResponse This struct is used to hold responses that are received from the server
type GameResponse struct {
	ID      int64
	Message string
}

func main() {
	db, _ := sql.Open("sqlite3", "dev.db")
	defer db.Close()

	DBStatus := db.Ping() == nil
	if DBStatus == false {
		fmt.Print("Could not connect to database")
	}

	http.HandleFunc("/startgame", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		word := RandomWordFromList()
		guess := CreateInitGuess(word)

		res, err := db.Exec("insert into games(id, user, word, guess) values (?, ?, ?, ?)",
			nil, user, word, guess)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		newID, err := res.LastInsertId()
		results := GameResponse{newID, "A new game has been started. Your hangman challenge is: " + guess}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		results := []Game{}

		rows, _ := db.Query("select * from games where user ='" + user + "'")
		for rows.Next() {
			var g Game
			rows.Scan(&g.ID, &g.User, &g.Word, &g.Guess)
			results = append(results, g)
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		gameid := r.FormValue("gameid")
		character := r.FormValue("character")
		results := GameResponse{}

		g := Game{}
		rows, _ := db.Query("select * from games where id =" + gameid + " and user ='" + user + "'")
		for rows.Next() {
			rows.Scan(&g.ID, &g.User, &g.Word, &g.Guess)
		}

		var ok bool
		g, ok = PlayGame(g, character)
		if true == ok {
			_, err := db.Exec("update games set guess = '" + g.Guess + "' where id = " + gameid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			msg := "Congratulations you uncovered a new letter. Your hangman challenge is: " + g.Guess
			if g.Guess == g.Word {
				msg = msg + ". Congratulations you have completed the game. "
			}
			results = GameResponse{g.ID, msg}
		} else {
			results = GameResponse{g.ID, "Unlucky that character does not exist. Your hangman challenge is: " + g.Guess}
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

// RandomWordFromList Function generates a random word as the hangman puzzle
func RandomWordFromList() string {
	words := [10]string{"computer", "house", "airplane", "helicopter", "table",
		"elephant", "giraffe", "software", "hardware", "france"}

	randIndex := rand.Intn(len(words))

	return words[randIndex]
}

// CreateInitGuess This function generates a blanked out string which is used as a placeholder when showing the hidden word to the user
func CreateInitGuess(word string) string {
	b := make([]byte, len(word))

	for i := range b {
		b[i] = '_'
	}

	return string(b)
}

// PlayGame This function replaces any matching characters that have been guessed correctly into the hidden string
func PlayGame(g Game, character string) (Game, bool) {
	found := false
	word := g.Word
	guess := g.Guess
	locations := []int{}

	locations, _ = FindCharacterLocations(word, character, locations, 0)
	if len(locations) > 0 {
		for _, value := range locations {
			guess = guess[:value] + character + guess[value+1:]
		}
	}

	if g.Guess != guess {
		found = true
	}
	g.Guess = guess

	return g, found
}

// FindCharacterLocations This is a recursive function to check multiple instances of a guessed character appearing in the hidden text
func FindCharacterLocations(word string, character string, locations []int, charCount int) ([]int, int) {
	i := strings.Index(word, character)

	if i > -1 {
		secondSplit := word[i+1:]
		locations = append(locations, charCount+i)

		if len(secondSplit) > 0 {
			charCount = charCount + i + 1
			locations, charCount = FindCharacterLocations(secondSplit, character, locations, charCount)
		}
	}

	return locations, charCount
}
