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

// Timeslot struct contains the timeslot schema in a struct format.
type Timeslot struct {
	ID        int    `json:"time_id"`
	UserID    int    `json:"user_id"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// GetTimeslot returns a timeslot from the database in JSON format, given the specific time_id as a URL parameter.
// This probably isn't needed at all.
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

// GetTimeslots returns all timeslots from the database for a specific user or company in JSON format.
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

// CreateTimeslot creates a new timeslot in the database. It expects user_id, day, start_time, and end_time.
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

// UpdateTimeslot updates a timeslot in the database, given its' time_id and other fields requesting to be changed.
func UpdateTimeslot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var timeslot Timeslot
	if err := json.NewDecoder(r.Body).Decode(&timeslot); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request parameters", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	timeslotID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE timeslots SET "

	structIterator := reflect.ValueOf(timeslot)
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
	sqlStatement = sqlStatement + " WHERE time_id" + "=$" + strconv.Itoa(j)
	values = append(values, timeslotID)
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

// DeleteTimeslot deletes a timeslot from the database given its' time_id. It will only work if
// the user sending the request has sufficient admin priveliges.
func DeleteTimeslot(w http.ResponseWriter, r *http.Request) {

}
