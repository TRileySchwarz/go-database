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

	"github.com/TRileySchwarz/go-database/auth"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/routes"
)

var TestUser = models.User{
	ID:        "Chad@gmail.com",
	Password:  "thisIsABadPassword1!",
	FirstName: "Chad",
	LastName:  "Chillerton",
}

// Verifies the login route is working as intended
// Does not test the unhappy path
func TestLogin(t *testing.T) {
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

	// Copy over the test data
	user := TestUser

	// hash the submitted password so its not stored in plain text
	hashedPass, err := auth.HashPassword(user.Password)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPass)

	// Add the user data to the DB
	_, err = routes.SignUpUser(&user)
	if err != nil {
		t.Fatalf("Could not signup user to DB: %v", err)
	}

	// Mock a login request to simulate attempt
	requestByte, _ := json.Marshal(models.LoginRequest{
		Email:    TestUser.ID,
		Password: TestUser.Password,
	})
	req, err := http.NewRequest("POST", "localhost:8080/login", bytes.NewReader(requestByte))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	routes.HandleLogin(rec, req)

	res := rec.Result()
	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Printf("\nThere is an error on body close: %v", err)
		}
	}()

	// Verify correct status code was received
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
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
		t.Errorf("expected the users email: %v; got %v", user.ID, email)
	}
}
