package routes

import (
	"encoding/json"
	"github.com/TRileySchwarz/go-database/models"
	"net/http"
)

// Helper used to set json response for API
// This just pushes all the error stack to the front end, ideally the front end would not be getting this much info
// Instead receive a series of predetermined error messages
func SetResponse(w http.ResponseWriter, errorMsg error, statusCode int) {

	// This only supports setting bad requests for now
	// This will become more useful when we start returning  predefined error codes indicating whats gone wrong
	// instead of just pushing the errors to front end
	if statusCode == http.StatusBadRequest {
		// Marshal the web token response
		responseJSON, err := json.Marshal(models.FrontEndErr{ErrorMsg: errorMsg.Error()})
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, err = w.Write(responseJSON)
		if err != nil {
			panic(err)
		}
	}
}