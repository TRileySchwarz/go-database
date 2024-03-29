package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/routes"
	"github.com/TRileySchwarz/go-database/auth"
)

// Verifies the /signup route is working as intended
// Does not test the unhappy path
func TestSignUp(t *testing.T) {
	// Load the environment variables, this is where things like api keys should be stored.
	// Can also store constants shared by multiple services
	err := godotenv.Load("../.env")
	if err != nil {
		panic(errors.Wrap(err, "Could not load .env file"))
	}

	// Connect to the database
	err = db.InitLocalDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database connection: %v", err)
	}

	// Mock a request
	requestByte, _ := json.Marshal(TestUser)
	req, err := http.NewRequest("POST", "localhost:8080/signup", bytes.NewReader(requestByte))
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

	// Verify correct status code was received
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK; got %v instead, with error: %v", res.StatusCode, err)
	}

	// Parse body of the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("there was an error reading the body of response: %v", err)
	}

	// Unmarshal the body so we can verify its contents
	var response models.WebTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("there was an error marshalling the body of the response: %v", err)
	}

	// Ensure the web token is valid
	email, err := auth.VerifyWebToken(response.Token)
	if err != nil {
		t.Fatalf("there was an error verifying the web token response: %v", err)
	}

	// Ensure the corresponding email address matches up
	if email != TestUser.ID {
		t.Errorf("expected the users email: %v; got %v", TestUser.ID, email)
	}
}
