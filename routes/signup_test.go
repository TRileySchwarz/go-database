package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/routes"
	"github.com/TRileySchwarz/go-database/webtoken"
)

func TestSignUp(t *testing.T) {
	// Connect to the database
	err := db.InitLocalDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database connection: %v", err)
	}

	user := models.User{
		ID:        "Chad@gmail.com",
		Password:  "thisIsABadPassword",
		FirstName: "Chad",
		LastName:  "Chillerton",
	}

	// TODO see if bytes.NewBuffer(requestbody) can replace this encoding function
	// Marshal the json from the get go
	// https://medium.com/@masnun/making-http-requests-in-golang-dd123379efe7
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(user)
	if err != nil {
		t.Fatalf("Could not encode user: %v", err)
	}

	req, err := http.NewRequest("POST", "localhost:8080/signup", b)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	routes.HandleSignUp(rec, req)

	res := rec.Result()
	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Printf("\nThere is an error on body close: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("there was an error reading the body of response: %v", err)
	}

	var response models.WebTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("there was an error marshalling the body of the response: %v", err)
	}

	email, err := webtoken.VerifyWebToken(response.Token)
	if err != nil {
		t.Fatalf("there was an error verifying the web token response: %v", err)
	}

	if email != user.ID {
		t.Errorf("expected the users email: %v; got %v", user.ID, email)
	}
}
