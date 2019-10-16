package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/webtoken"
)

// HandleUsers switches on the GET and PUT requests to provide users access
// to changing their information, or returning all user data
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	// Pull out the email to the corresponding web token
	webTokenString := r.Header.Get("x-authentication-token")
	email, err := webtoken.VerifyWebToken(webTokenString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Switch on the get and put methods to determine which route is being called
	if r.Method == http.MethodGet {
		handleGetUsers(w)
	} else if r.Method == http.MethodPut {
		handlePutUsers(w, r, email)
	} else {
		// Signal a bad request was sent, ie. not a get or put
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Responsible for handling the GET version of /users
func handleGetUsers(w http.ResponseWriter) {
	u, err := getUsers()
	if err != nil {
		panic(err)
	}

	// Remove all the passwords from the struct before returning it
	resp := removePasswords(u)

	// Marshal the response before returning it
	respJSON, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respJSON)
	if err != nil {
		panic(err)
	}
}

// Returns all of the user data stored in the DB
func getUsers() ([]models.User, error) {
	var u []models.User
	err := db.DataBase.Model(&u).Select()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Helper function to remove the passwords from the struct slice before returning it
func removePasswords(users []models.User) []models.UserNoPwd {
	var usersToReturn []models.UserNoPwd

	for k := range users {
		userToAppend := models.UserNoPwd{
			Email:     users[k].ID,
			FirstName: users[k].FirstName,
			LastName:  users[k].LastName,
		}

		usersToReturn = append(usersToReturn, userToAppend)
	}

	return usersToReturn
}

// Responsible for handling the PUT version of /users
func handlePutUsers(w http.ResponseWriter, r *http.Request, email string) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the request body into our User struct
	var newUserInfo models.PutUserRequest
	err = json.Unmarshal(body, &newUserInfo)
	if err != nil {
		panic(err)
	}

	// Update the user info in the database
	err = updateUser(email, &newUserInfo)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Handles the request for a user to update their first and last names
func updateUser(email string, request *models.PutUserRequest) error {
	// Select user by primary key.
	user := &models.User{ID: email}
	err := db.DataBase.Select(user)
	if err != nil {
		return err
	}

	// Update fields
	user.FirstName = request.FirstName
	user.LastName = request.LastName

	// Apply the updated user struct to the database
	err = db.DataBase.Update(user)
	if err != nil {
		return err
	}

	return nil
}
