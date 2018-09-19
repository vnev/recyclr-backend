package handlers

import (
	"encoding/json"
	"net/http"
)

// GetUsers : function to return a user from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []string // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
