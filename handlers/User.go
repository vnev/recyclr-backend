package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User : basic user schema
type User struct {
	ID   string
	Name string
}

// GetUsers : function to return a user from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser : function to return a user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var users []User // TODO: actually get this to read in users from the DB
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

// CreateUser : function to create a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	var user User // actually need to create User struct
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID - not safe
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}
