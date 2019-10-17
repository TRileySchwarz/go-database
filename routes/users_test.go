package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TRileySchwarz/go-database/db"
	"github.com/TRileySchwarz/go-database/models"
	"github.com/TRileySchwarz/go-database/routes"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var TestUsers = []models.User{{
		ID:        "Piccolo@gmail.com",
		Password:  "thisIsABadPassword",
		FirstName: "Piccolo",
		LastName:  "1",
	}, {
		ID:        "Goku@gmail.com",
		Password:  "thisIsAHorriblePassword",
		FirstName: "Goku",
		LastName:  "2",
	}, {
		ID:        "Frieza@gmail.com",
		Password:  "thisIsATerriblePassword",
		FirstName: "Frieza",
		LastName:  "3",
	},
}

var TestUsersNoPwd = []models.UserNoPwd{{
		Email:        "Piccolo@gmail.com",
		FirstName: "Piccolo",
		LastName:  "1",
	}, {
		Email:        "Goku@gmail.com",
		FirstName: "Goku",
		LastName:  "2",
	}, {
		Email:        "Frieza@gmail.com",
		FirstName: "Frieza",
		LastName:  "3",
	},
}

// Verifies the GET /users route is working as intended
// Does not test the unhappy path
func TestUsersGet(t *testing.T) {
	// Load the environment variables, this is where things like api keys should be stored.
	// Can also store constants shared by multiple services
	err := godotenv.Load("../.env")
	if err != nil {
		panic(errors.Wrap(err, "Could not load .env file"))
	}

	// Connect/Initialize to the database
	err = db.InitLocalDatabase()
	if err != nil {
		t.Fatalf("Could not initialize database connection: %v", err)
	}

	// Populate the database with test data
	jwts := populateDatabase(t)

	// Iterate through the tokens and get the user database ensuring that all jwts work as intended
	for _, token := range jwts {
		// Mock a request
		req, err := http.NewRequest("GET", "localhost:8080/users", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		req.Header.Set("x-authentication-token", token.Token)

		rec := httptest.NewRecorder()
		routes.HandleUsers(rec, req)

		res := rec.Result()
		// TODO Check defer resource leaks in for loop
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
		var response []models.UserNoPwd
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatalf("there was an error marshalling the body of the response: %v", err)
		}

		// Ensure the response matches up as intended
		if !reflect.DeepEqual(response, TestUsersNoPwd) {
			t.Errorf("expected response %v: got %v\n", TestUsersNoPwd, response)
		}
	}
}

// Verifies the PUT /users route is working as intended
// Does not test the unhappy path
func TestUsersPut(t *testing.T) {
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

	// Copy over test user
	user := TestUser

	// Signup user and return a valid JWT
	jwt, err := routes.SignUpUser(&user)
	if err != nil {
		t.Fatalf("Could not signup user to DB: %v", err)
	}

	// Swap the first and last names
	requestByte, _ := json.Marshal(models.PutUserRequest{
		FirstName: TestUser.LastName,
		LastName:  TestUser.FirstName,
	})

	// Create new request
	req, err := http.NewRequest("PUT", "localhost:8080/users", bytes.NewReader(requestByte))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("x-authentication-token", jwt.Token)

	rec := httptest.NewRecorder()
	routes.HandleUsers(rec, req)

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

	// Now we need to check the database updated the values correctly
	req, err = http.NewRequest("GET", "localhost:8080/users", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("x-authentication-token", jwt.Token)

	rec = httptest.NewRecorder()
	routes.HandleUsers(rec, req)
	res = rec.Result()

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
	var response []models.UserNoPwd
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("there was an error marshalling the body of the response: %v", err)
	}

	// The expected result with the first and last name switched
	userExpected := []models.UserNoPwd{{
			Email:     TestUser.ID,
			FirstName: TestUser.LastName,
			LastName:  TestUser.FirstName,
		},
	}

	// Ensure the corresponding email address matches up
	if !reflect.DeepEqual(response, userExpected) {
		t.Errorf("expected response %v: got %v\n", userExpected, response)
	}
}

// Populates the database with a series of user test data, returns the corresponding slice of jwts
func populateDatabase(t *testing.T) []models.WebTokenResponse {
	var jwts []models.WebTokenResponse

	// Range over user data to add users to database
	for _, u := range TestUsers {
		jwt, err := routes.SignUpUser(&u)
		if err != nil {
			t.Fatalf("Could not signup user to DB: %v", err)
		}

		// add the corresponding jwts to a list to check
		jwts = append(jwts, jwt)
	}

	return jwts
}
