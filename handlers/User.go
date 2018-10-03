package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	Password  string `json:"passwd"`
	Token     string `json:"token"`
}

// GetUser : function to return a user from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	// this returns a blank password field but it looks jank anyway
	// TODO: maybe just return a newly defined struct without password field
	var user User
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT user_id, address, email, user_name, is_company, rating, joined_on FROM users WHERE user_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Address, &user.Email, &user.Name, &user.IsCompany, &user.Rating, &user.JoinedOn)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		panic(err)
	}

	json.NewEncoder(w).Encode(&user)
}

// CreateUser : function to create a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO users (address, email, user_name, is_company, passwd)
	VALUES ($1, $2, $3, $4, crypt($5, gen_salt('md5')))
	RETURNING user_id`
	id := 0
	err := db.DBconn.QueryRow(sqlStatement, user.Address, user.Email, user.Name, false, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// fmt.Println("New user created with ID:", id)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser : function to update a user in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.ID == 0 {
		http.Error(w, "No user ID found", http.StatusBadRequest)
	}

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE users SET "

	// Time to iteratively loop over a struct, the easiest to understand syntax ever!
	structIterator := reflect.ValueOf(user)
	for i := 0; i < structIterator.NumField(); i++ {
		//fmt.Printf("field: %+v, value: %+v\n", structIterator.Type().Field(i).Name, structIterator.Field(i).Interface())
		field := structIterator.Type().Field(i).Name
		val := structIterator.Field(i).Interface()

		// Check if the field is zero-valued, meaning it won't be updated
		//fmt.Printf("VAL IS %v and TYPE IS %v and ZERO OF TYPE IS %v\n", val, structIterator.Field(i).Type(), reflect.Zero(structIterator.Field(i).Type()).Interface())
		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			// fmt.Printf("%v is non-zero, adding to update\n", field)
			sqlStatement = sqlStatement + strings.ToLower(field) + "=$" + strconv.Itoa(j) + ", "
			j++
			values = append(values, val)
		}
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-2]
	sqlStatement = sqlStatement + " WHERE user_id" + "=$" + strconv.Itoa(j)
	values = append(values, userID)
	// fmt.Printf("executing SQL: \n\t%s\n", sqlStatement)
	// fmt.Printf("$1 is %s and $2 is %d\n", values[0], values[1])
	row, err := db.DBconn.Exec(sqlStatement, values...) //.Scan(&user.ID, &user.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	count, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(count, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	//json.NewEncoder(w).Encode({"status": "200", "message": "success"})
}

// AuthenticateUser : generate JWT for user and return
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	sqlStatement := "SELECT email FROM users WHERE email=$1"
	email := ""
	_ = db.DBconn.QueryRow(sqlStatement, user.Email).Scan(&email)
	if email == "" {
		http.Error(w, "No user found", http.StatusBadRequest)
		return
	}

	userID := 0
	sqlStatement = "SELECT user_id FROM users WHERE email=$1 AND passwd=crypt($2, passwd)"
	err := db.DBconn.QueryRow(sqlStatement, user.Email, user.Password).Scan(&userID)
	if err != nil {
		http.Error(w, "No user found with that email/password", http.StatusBadRequest)
		return
	}
	// if err != nil {
	// 	http.Error(w, err.Error(), HTTPInternalError)
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "recyclr.xyz",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"name": user.Name,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement = "UPDATE users SET token=$1 WHERE user_id=$2"
	_, err = db.DBconn.Exec(sqlStatement, tokenString, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["token"] = tokenString

	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// LogoutUser : This logs a user out
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.ID == 0 {
		http.Error(w, "No user ID found", http.StatusBadRequest)
		return
	}

	sqlStatement := "UPDATE users SET token='0' WHERE user_id=$1"
	_, err := db.DBconn.Exec(sqlStatement, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteUser : function to delete a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// var user User
	// _ = json.NewDecoder(r.Body).Decode(&user)
	// if user.ID == 0 {
	// 	http.Error(w, "No user ID found", http.StatusBadRequest)
	// }

	// sqlStatement :=
}

// GetProgress : function to get the progress of a user's listings
func GetProgress(w http.ResponseWriter, r *http.Request) {
	var listings []Listing
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rows, err := db.DBconn.Query("SELECT listing_id, title, description, material_type, material_weight, active FROM listings WHERE user_id=$1", userID)
	//err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Address, &user.Email, &user.Name, &user.IsCompany, &user.Rating, &user.JoinedOn)

	defer rows.Close()
	for rows.Next() {
		var listing Listing
		err = rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.MaterialType, &listing.MaterialWeight, &listing.Active)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		listings = append(listings, listing)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(listings)
}
