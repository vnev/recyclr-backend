package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

// User : basic user schema
type User struct {
	ID        int    `json:"user_id"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IsCompany bool   `json:"is_company"`
	Rating    int    `json:"rating"`
	JoinedOn  string `json:"joined_on"`
	Password  string `json:"password"`
}

// GetUser : function to return a user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT user_id, user_name FROM users WHERE user_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 500)
			return
		}
		panic(err)
	}

	json.NewEncoder(w).Encode(&user)
}

// GetUsers : function to return a user from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DBconn.Query("SELECT user_id, user_name FROM users")
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		fmt.Printf("ID is %d, Name is %s\n", user.ID, user.Name)
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(users)
}

// CreateUser : function to create a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO users (address, email, user_name, is_company, password)
	VALUES ($1, $2, $3, $4, crypt($5, gen_salt('md5')))
	RETURNING user_id`
	id := 0
	err := db.DBconn.QueryRow(sqlStatement, user.Address, user.Email, user.Name, false, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("New user created with ID:", id)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser : function to update a user in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	/*var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var newUser User
			_ = json.NewDecoder(r.Body).Decode(&newUser)
			newUser.ID = params["id"]
			users = append(users, newUser)
			json.NewEncoder(w).Encode(newUser)
			return
		}
	}
	json.NewEncoder(w).Encode(users)*/
}

// DeleteUser : function to delete a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	/*var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)*/
}
