package handlers

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/vnev/recyclr-backend/db"
)

func TestGetTimeslots(t *testing.T) {
	var timeslots []Timeslot
	//w.Header().Set("Content-Type", "application/json")

	//params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi("1")
	if err != nil {
		t.Error(err.Error())
	}

	rows, err := db.DBconn.Query("SELECT time_id, day, start_time, end_time FROM timeslots WHERE user_id=$1", userID)
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err.Error())
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
		t.Error(err.Error())
	}
}
