package main

import (
	"fmt"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize the database and create corresponding rows
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}

	// Defer closing the database when the program exits
	defer func() {
		err = db.DataBase.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Assign the router handlers for this particular API
	http.HandleFunc("/signup", routes.HandleSignUp)
	http.HandleFunc("/login", routes.HandleLogin)
	http.HandleFunc("/users", routes.HandleUsers)

	// Initialize the API and prepare to handle requests
	fmt.Println("Starting the API Service...")
	fmt.Println("API is being served on port: 8093")
	err = http.ListenAndServe(":8093", nil)
	if err != nil {
		log.Print(err)
	}
}
