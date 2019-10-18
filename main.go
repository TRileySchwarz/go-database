package main

import (
	"fmt"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/routes"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func main() {
	// Load the environment variables, this is where things like api keys should be stored.
	// Can also store constants shared by multiple services
	err := godotenv.Load()
	if err != nil {
		panic(errors.Wrap(err, "Could not load .env file"))
	}

	apiPort := os.Getenv("API_PORT")

	// Initialize the database and create corresponding rows
	err = db.InitDatabase()
	if err != nil {
		panic(errors.Wrap(err, "There was an error initializing the database connection"))
	}
	fmt.Println("\nConnection to the database successful")

	// Defer closing the database when the program exits
	defer func() {
		err = db.DataBase.Close()
		if err != nil {
			fmt.Printf("\nError closing the db connection: %v", err)
		}
	}()

	// Initialize the API and prepare to handle requests
	fmt.Println("Starting the API Service on port: " + apiPort)
	err = http.ListenAndServe(":" + apiPort, handler())
	if err != nil {
		panic(errors.Wrap(err, "Could not listen and serve on port: " + apiPort))
	}
}

// http handler we will be using to route API calls
func handler() http.Handler {
	r := http.NewServeMux()

	// Assign the router handlers for this particular API
	r.HandleFunc("/signup", routes.HandleSignUp)
	r.HandleFunc("/login", routes.HandleLogin)
	r.HandleFunc("/users", routes.HandleUsers)

	return r
}
