package main

import (
	"fmt"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {

	// Load the environment variables, this is where things like api keys should be stored.
	// Can also store constants shared by multiple services
	godotenv.Load()

	// Initialize the database and create corresponding rows
	err := db.InitDatabase()
	if err != nil {
		fmt.Printf("There is db error: %v", err)
	}

	fmt.Println("Connection to the database successful")

	// Defer closing the database when the program exits
	defer func() {
		err = db.DataBase.Close()
		if err != nil {
			fmt.Printf("There is db error on close: %v", err)
		}
	}()

	// Assign the router handlers for this particular API
	http.HandleFunc("/signup", routes.HandleSignUp)
	http.HandleFunc("/login", routes.HandleLogin)
	http.HandleFunc("/users", routes.HandleUsers)

	// Initialize the API and prepare to handle requests
	fmt.Println("Starting the API Service...")
	fmt.Println("API is being served on port: 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
}
