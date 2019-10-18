package routes

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/auth"
)

// HandleUsers switches on the GET and PUT requests to provide users access
// to changing their information, or returning all user data
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	// Pull out the email to the corresponding web token
	webTokenString := r.Header.Get("x-authentication-token")
	email, err := auth.VerifyWebToken(webTokenString)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Switch on the get and put methods to determine which route is being called
	if r.Method == http.MethodGet {
		handleGetUsers(w)
	} else if r.Method == http.MethodPut {
		handlePutUsers(w, r, email)
	} else {
		// Signal a bad request was sent, ie. not a get or put
		SetResponse(w, errors.New("invalid method used"), http.StatusBadRequest)
	}
}

// Responsible for handling the GET version of /users
func handleGetUsers(w http.ResponseWriter) {
	u, err := getUsers()
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Marshal the response before returning it
	respJSON, err := json.Marshal(u)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respJSON)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}
}

// Returns all of the user data stored in the DB
func getUsers() ([]models.User, error) {
	var u []models.User
	err := db.DataBase.Model(&u).Select()
	if err != nil {
		return nil, errors.New("issue retrieving all the users from db")
	}

	return u, nil
}

// Responsible for handling the PUT version of /users
func handlePutUsers(w http.ResponseWriter, r *http.Request, email string) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into our User struct
	var newUserInfo models.PutUserRequest
	err = json.Unmarshal(body, &newUserInfo)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Update the user info in the database
	err = updateUser(email, &newUserInfo)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
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
		return errors.Wrap(err, "issue locating the user information")
	}

	// Update fields
	user.FirstName = request.FirstName
	user.LastName = request.LastName

	// Apply the updated user struct to the database
	err = db.DataBase.Update(user)
	if err != nil {
		return errors.Wrap(err, "issue updating the users info")
	}

	return nil
}
