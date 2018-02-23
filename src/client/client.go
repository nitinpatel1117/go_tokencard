package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

// Game This struct holds all the game data
type Game struct {
	ID    int64  `json:"ID"`
	User  string `json:"User"`
	Word  string `json:"Word"`
	Guess string `json:"Guess"`
}

// GameResponse This struct is used to hold responses that are received from the server
type GameResponse struct {
	ID      int64  `json:"ID"`
	Message string `json:"Message"`
}

func main() {
	var user string
	var gameid string
	var char string

	app := cli.NewApp()
	app.Name = "Hangman"
	app.Usage = "Play the game by guessing the missing letters"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"c"},
			Usage:   "start a new game",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "user, u",
					Usage:       "the name of the user who wants to start a new game",
					Destination: &user,
				},
			},
			Action: func(c *cli.Context) error {
				g, err := startgame(user)
				failOnError(err, "Unable to start a new game for "+user)

				fmt.Println("\n" + g.Message)
				fmt.Println("The id of your game is: " + strconv.FormatInt(g.ID, 10))
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"a"},
			Usage:   "list all games that are currently playing",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "user, u",
					Usage:       "the name of the user who's games we want to list",
					Destination: &user,
				},
			},
			Action: func(c *cli.Context) error {
				games, err := list(user)
				failOnError(err, "Unable to retrieve user's list of games")

				displayGames(games)
				return nil
			},
		},
		{
			Name:    "guess",
			Aliases: []string{"a"},
			Usage:   "play a game",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "user, u",
					Usage:       "the name of the user who is playing the game",
					Destination: &user,
				},
				cli.StringFlag{
					Name:        "gameid, g",
					Usage:       "the name of the user who is playing the game",
					Destination: &gameid,
				},
				cli.StringFlag{
					Name:        "char, c",
					Usage:       "the name of the user who is playing the game",
					Destination: &char,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("guess", user, gameid, char)

				g, err := guess(user, gameid, char)
				failOnError(err, "Unable to submit the character guess for user: "+user)

				fmt.Println("\n" + g.Message)
				fmt.Println("The id of your game is: " + strconv.FormatInt(g.ID, 10))
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func startgame(user string) (GameResponse, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get("http://localhost:8080/startgame?user=" + user); err != nil {
		return GameResponse{}, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return GameResponse{}, err
	}

	var data GameResponse
	err = json.Unmarshal(body, &data)

	return data, err
}

func list(user string) ([]Game, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get("http://localhost:8080/list?user=" + user); err != nil {
		return []Game{}, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return []Game{}, err
	}

	var data []Game
	err = json.Unmarshal(body, &data)

	return data, err
}

func guess(user string, gameid string, char string) (GameResponse, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get("http://localhost:8080/guess?user=" + user + "&gameid=" + gameid + "&character=" + char + ""); err != nil {
		return GameResponse{}, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return GameResponse{}, err
	}

	var data GameResponse
	err = json.Unmarshal(body, &data)

	return data, err
}

func displayGames(games []Game) {
	fmt.Println("| gameid | user | guess | complete |")
	for _, g := range games {

		var completed string
		completed = "No"
		if g.Guess == g.Word {
			completed = "Yes"
		}

		fmt.Println("| " + strconv.FormatInt(g.ID, 10) + " | " + g.User + " | " + g.Guess + " | " + completed + " |")
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
