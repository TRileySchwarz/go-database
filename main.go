package main

import (
	"fmt"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/routes"
	"log"
	"net/http"
)

func main() {
	//// Initialize the database and create corresponding rows
	//err := db.InitDatabase()
	//if err != nil {
	//	fmt.Printf("There is db error: %v", err)
	//}
	//
	//// Defer closing the database when the program exits
	//defer func() {
	//	err = db.DataBase.Close()
	//	if err != nil {
	//		fmt.Printf("There is db error on close: %v", err)
	//	}
	//}()

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})


	// Assign the router handlers for this particular API
	http.HandleFunc("/signup", routes.HandleSignUp)
	http.HandleFunc("/login", routes.HandleLogin)
	http.HandleFunc("/users", routes.HandleUsers)

	// TODO hack way of connecting to database
	http.HandleFunc("/db", HandleDBInit)

	// Initialize the API and prepare to handle requests
	fmt.Println("Starting the API Service...")
	fmt.Println("API is being served on port: 8093")
	err := http.ListenAndServe(":8093", nil)
	if err != nil {
		log.Print(err)
	}
}

// Handles the users route and switched on the
func HandleDBInit(w http.ResponseWriter, r *http.Request) {

	// Initialize the database and create corresponding rows
	err := db.InitDatabase()
	if err != nil {
		fmt.Fprintf(w, "There is db error: %v", err)
	}

	// Defer closing the database when the program exits
	defer func() {
		err = db.DataBase.Close()
		if err != nil {
			fmt.Fprintf(w, "There is db error on close: %v", err)
		}
	}()
}

