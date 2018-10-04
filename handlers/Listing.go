package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

// Listing : basic listing schema
type Listing struct {
	ID             int     `json:"listing_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ImageHash      string  `json:"img_hash"`
	MaterialType   string  `json:"material_type"`
	MaterialWeight float64 `json:"material_weight"`
	UserID         int     `json:"user_id"`
	Active         bool    `json:"is_active"`
	PickupDate     string  `json:"pickup_date"`
}

// GetListing : function to return a listing from the database
func GetListing(w http.ResponseWriter, r *http.Request) {
	var listing Listing
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT title, description, img_hash, material_type, material_weight, active FROM listings WHERE listing_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, listingID).Scan(&listing.Title, &listing.Description, &listing.ImageHash, &listing.MaterialType, &listing.MaterialWeight, &listing.Active)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&listing)
}

// GetListings : function to return all listings from the database
func GetListings(w http.ResponseWriter, r *http.Request) {
	var listings []Listing
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DBconn.Query("SELECT listing_id, title, description, material_type, material_weight, active FROM listings")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var listing Listing
		err = rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.MaterialType, &listing.MaterialWeight, &listing.Active)
		fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		listings = append(listings, listing)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listings)
}

// CreateListing : function to create a new listing in the database
func CreateListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var listing Listing
	err := json.NewDecoder(r.Body).Decode(&listing)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}
	fmt.Println("LISTING IS: ", listing)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO listings (title, description, img_hash, material_type, material_weight, user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING listing_id`
	id := 0
	err = db.DBconn.QueryRow(sqlStatement, listing.Title, listing.Description, listing.ImageHash, listing.MaterialType, listing.MaterialWeight, listing.UserID).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("New listing created with ID: ", id)
	json.NewEncoder(w).Encode(listing)
}

// UpdateListing : function to update a listing in the database
func UpdateListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var listing Listing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request parameters", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE listings SET "

	structIterator := reflect.ValueOf(listing)
	for i := 0; i < structIterator.NumField(); i++ {
		field := structIterator.Type().Field(i).Name
		val := structIterator.Field(i).Interface()

		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			sqlStatement += strings.ToLower(field) + "=$" + strconv.Itoa(j) + ", "
			j++
			values = append(values, val)
		}
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-2]
	sqlStatement = sqlStatement + " WHERE listing_id" + "=$" + strconv.Itoa(j)
	values = append(values, listingID)
	row, err := db.DBconn.Exec(sqlStatement, values...)
	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteListing : function to delete a listing from the database
func DeleteListing(w http.ResponseWriter, r *http.Request) {
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
