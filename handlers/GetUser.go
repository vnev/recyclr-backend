package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetUser : function to return a user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var users []string // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through users and find with id
	for _, user := range users {
		if user.ID == params["id"] { // needs to actually be set to User struct
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{}) // also will be fixed once user struct exists
}
