package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

// Timeslot : struct to hold timeslot information
type Timeslot struct {
	ID        int    `json:"time_id"`
	UserID    int    `json:"user_id"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// GetTimeslot : function to return a timeslot from the database, probably unneeded
func GetTimeslot(w http.ResponseWriter, r *http.Request) {
	var timeslot Timeslot
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	timeID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT user_id, day, start_time, end_time FROM timeslots WHERE time_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, timeID).Scan(&timeslot.UserID, &timeslot.Day, &timeslot.StartTime, &timeslot.EndTime)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&timeslot)
}

// GetTimeslots : function to return all timeslots for a company/user from the database
func GetTimeslots(w http.ResponseWriter, r *http.Request) {
	var timeslots []Timeslot
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := db.DBconn.Query("SELECT time_id, day, start_time, end_time FROM timeslots WHERE user_id=$1", userID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var timeslot Timeslot
		err = rows.Scan(&timeslot.ID, &timeslot.Day, &timeslot.StartTime, &timeslot.EndTime)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		timeslots = append(timeslots, timeslot)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(timeslots)
}

// CreateTimeslot : function to create a new listing in the database
func CreateTimeslot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var timeslot Timeslot
	err := json.NewDecoder(r.Body).Decode(&timeslot)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}
	//fmt.Println("LISTING IS: ", listing)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO timeslots (user_id, day, start_time, end_time)
	VALUES ($1, $2, $3, $4)
	RETURNING time_id`
	id := 0
	err = db.DBconn.QueryRow(sqlStatement, timeslot.UserID, timeslot.Day, timeslot.StartTime, timeslot.EndTime).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("New listing created with ID: ", id)
	json.NewEncoder(w).Encode(timeslot)
}
