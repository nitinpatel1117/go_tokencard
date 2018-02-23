# exercise-go-dev

This is a small server and client that plays a hangman game

## Setup

This game makes use of a backend database to store the game data. It has been setup to use sqlite3 just to get going quickly.

### Setup database
You can download and setup sqlite3 from https://www.sqlite.org/quickstart.html 
Once you have sqlite3 running on your machine, please navigate to the  directory `./src/server` and run `sqlite3 dev.db` to confirm that you can connect to the database. once you are in the database you may exit by typing `.exit`

### Libraries
 The application uses the below libraries
`go get github.com/stretchr/testify/assert`
`go get github.com/urfave/cli`
`go get github.com/mattn/go-sqlite3`
`go get github.com/stretchr/testify`


## Start the server
The server can be started by navigating to `src/server` and then issue the following command `go run server.go`

## Playing the game
The game is played via the client. Navigate to `src/client`
and then the following commands can be issued to play the game

#### Start a new game
Issue the command:
`go run client.go start --user joe`
This will start a new game for the user joe. The command will return the id of the game as well as the current status of the hangman game

#### List all games that a user is playing
Issue the command:
`go run client.go list --user joe`
This will list all the games that the user joe is playing as well as the progres of that game and whether that game is completed or not

#### Play a game
To play a game you will need the game id which is returned by both the above calls to start and list a game,
Issue the command:
`go run client.go guess --user joe --gameid 29 --char o`

where the arguments:
 - user =  is the user who is playing the game
 - gameid =  is the id of the game
 - char =  is the letter that is being guessed

The command will return the result of whether a letter was found or not as well as the word puzzle with missing letters replaced with '_'.


## Running the test
There is a test file in `src/server/server_test.go`
This can be run with the command `go test`


## Additional things to do
These are additional things that I would have done if I had more time:
 1. Refactored the server.go  and client.go files out more. e.g.
		 - The server.go has all the database querying which can be put into another file in the `src/server`.
		 - The structs are copied in both the server.go and client.go file, they could go in another shared file.
		 - The url construction in the client.go can be separated
		 - The http client aspect of the client.go that calls the urls is similar in all three cases, some more refactoring here to not repeat the code
2. Add validation on the user input to the cli arguments
3. Add logging
