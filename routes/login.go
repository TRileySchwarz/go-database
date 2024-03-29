package routes

import (
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/TRileySchwarz/go-database/auth"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"

	"io/ioutil"
	"net/http"
)

// HandleLogin provides users the ability to submit a login request and returns a JWT if succesful
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Unmarshal the request details
	var loginDetails models.LoginRequest
	err = json.Unmarshal(body, &loginDetails)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Verify the login details and return a new json web token string
	// Should not crash the system if the login fails, hence no panic(err)
	jwt, err := verifyLoginDetails(loginDetails)
	if err != nil {
		SetResponse(w, err, http.StatusBadRequest)
		return
	}

	// Marshal teh web token response
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

// Verifies the login details and returns a jwt string in response
func verifyLoginDetails(loginDetails models.LoginRequest) (models.WebTokenResponse, error) {
	// Select user by primary key.
	user := &models.User{ID: loginDetails.Email}
	err := db.DataBase.Select(user)
	if err != nil {
		return models.WebTokenResponse{}, errors.Wrap(err, "invalid user email")

	}

	// Compare the hashed password with the one stored
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password))
	if err != nil {
		return models.WebTokenResponse{}, errors.Wrap(err, "invalid login credentials")
	}

	// Generates a web token string
	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return models.WebTokenResponse{}, errors.Wrap(err, "failed to generate web token")
	}

	return models.WebTokenResponse{Token: tokenString}, nil
}
