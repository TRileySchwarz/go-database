package routes

import (
	"encoding/json"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/webToken"
	"io/ioutil"
	"net/http"
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the request body into our User struct
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	// Pass the new User data to the database and attempt to insert it
	jwt, err := signUpUser(&user)
	if err != nil {
		panic(err)
	}

	// Marshal the web token response
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

// Helper used to signup a user in the DB, returns the corresponding JWT with those details
func signUpUser(user *models.User) (models.WebTokenResponse, error) {
	err := db.DataBase.Insert(user)
	if err != nil {
		panic(err)
	}

	tokenString, err := webToken.GenerateJWT(user.ID)
	if err != nil {
		return models.WebTokenResponse{}, err
	}

	return models.WebTokenResponse{tokenString}, nil
}
