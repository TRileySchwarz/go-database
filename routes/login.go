package routes

import (
	"encoding/json"
	"errors"

	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/webtoken"

	"io/ioutil"
	"net/http"
)

// HandleLogin provides users the ability to submit a login request and returns a JWT if succesful
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the request details
	var loginDetails models.LoginRequest
	err = json.Unmarshal(body, &loginDetails)
	if err != nil {
		panic(err)
	}

	// Verify the login details and return a new json web token string
	// Shouldnt crash the system if the login fails, hence no panic(err)
	jwt, err := verifyLoginDetails(loginDetails)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Marshal teh web token response
	responseJSON, err := json.Marshal(jwt)
	if err != nil {
		panic(err)
	}

	// Set the response header and write the body payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		panic(err)
	}
}

// Verifies the login details and returns a jwt string in response
func verifyLoginDetails(loginDetails models.LoginRequest) (models.WebTokenResponse, error) {
	// Select user by primary key.
	user := &models.User{ID: loginDetails.Email}
	err := db.DataBase.Select(user)
	if err != nil {
		return models.WebTokenResponse{}, err
	}

	// Check that the supplied password matches the one we have stored in the DB
	if loginDetails.Password != user.Password {
		return models.WebTokenResponse{}, errors.New("invalid login credentials")
	}

	// Generates a web token string
	tokenString, err := webtoken.GenerateJWT(user.ID)
	if err != nil {
		return models.WebTokenResponse{}, err
	}

	return models.WebTokenResponse{Token: tokenString}, nil
}
