package handlers

import (
	"testing"

	"github.com/vnev/recyclr-backend/db"
)

func TestGetMessages(t *testing.T) {
	var messages []Message
	type getAttributes struct {
		ForListing int `json:"for_listing"`
	}

	var attr getAttributes
	attr.ForListing = 1
	//w.Header().Set("Content-Type", "application/json")

	//_ = json.NewDecoder(r.Body).Decode(&attr)

	sqlStatement := `SELECT message_id, from_user, to_user, message_time, message_content FROM Messages WHERE for_listing=$1 ORDER BY "message_time"`
	rows, err := db.DBconn.Query(sqlStatement, attr.ForListing)
	if err != nil {
		t.Error(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var message Message
		if err = rows.Scan(&message.ID, &message.FromUser, &message.ToUser, &message.Timestamp, &message.Content); err != nil {
			t.Error(err.Error())
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		t.Error(err.Error())
	}
}
