package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

// CreateUser : function to create a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var users []string // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	var user User // actually need to create User struct
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID - not safe
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}
