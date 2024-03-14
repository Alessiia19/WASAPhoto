package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin handles the login request.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	//  Set the HTTP response header to indicate that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")

	//  Extract the username from the request body.
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		//  If there is an error decoding the username, returns a Bad Request status.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Login: Invalid request")
		return
	}

	// Check if the username meets the requirements
	if !isValidUsername(user.Username) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("Login: Invalid username format. Please follow the specified requirements.")
		return
	}

	//  Attempt to create a new user or retrieve an existing user from the database.
	newUser, err := rt.db.CreateUser(user.UserToDatabase())
	if err != nil {
		//  If there is an error creating or retrieving the user, returns an Internal Server Error status.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  Update the user data with the information from the database.
	user.UserFromDatabase(newUser)

	//  Returns a Created status and encode the user data in the response body.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
