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

// HandleSignUp allows a new user to register for the system. 
// The email can only be used once ie. unique
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into our User struct
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Verify the passwords strength
	pass := auth.ProcessPassword(user.Password)
	if pass.Score  <  auth.MinPassStrength {
		SetResponse(w, errors.New("Password is not strong enough"), http.StatusBadRequest)
		return
	}

	// hash the submitted password so its not stored in plain text
	hashedPass, err := auth.HashPassword(user.Password)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}
	user.Password = string(hashedPass)

	// Pass the new User data to the database and attempt to insert it
	jwt, err := SignUpUser(&user)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Marshal the web token response
	responseJSON, err := json.Marshal(jwt)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Set the response header and write the body payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
	}
}

// SignUpUser is a helper used to signup a user in the DB, 
// returns the corresponding JWT with those details
func SignUpUser(user *models.User) (models.WebTokenResponse, error) {
	err := db.DataBase.Insert(user)
	if err != nil {
		return models.WebTokenResponse{}, errors.Wrap(err, "Issue inserting the user into db")
	}

	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return models.WebTokenResponse{}, errors.Wrap(err, "Issue generating the jwt")
	}

	return models.WebTokenResponse{tokenString}, nil
}
