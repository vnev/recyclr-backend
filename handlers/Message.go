package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vnev/recyclr-backend/db"
)

// Message struct to store message information
type Message struct {
	ID        int    `json:"message_id"`
	Timestamp string `json:"message_time"`
	FromUser  int    `json:"from_user"`
	ToUser    int    `json:"to_user"`
	ListingID int    `json:"for_listing"`
	Content   string `json:"message_content"`
}

// GetMessages returns all messages between user and company for a particular listing
func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	type getAttributes struct {
		FromUser   int `json:"from_user"`
		ToUser     int `json:"to_user"`
		ForListing int `json:"for_listing"`
	}

	var attr getAttributes
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewDecoder(r.Body).Decode(&attr)

	sqlStatement := `SELECT message_id, message_time, message_content FROM Messages WHERE from_user=$1 AND to_user=$2 AND for_listing=$3 ORDER BY "message_time"`
	rows, err := db.DBconn.Query(sqlStatement, attr.FromUser, attr.ToUser, attr.ForListing)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var message Message
		if err = rows.Scan(&message.ID, &message.Timestamp, &message.Content); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// PutMessage adds a new message between user and company for particular listing
func PutMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage Message
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewDecoder(r.Body).Decode(&newMessage)

	newMessageID := -1
	sqlStatement := "INSERT INTO Messages (message_time, from_user, to_user, for_listing, message_content) VALUES ($1,$2,$3,$4,$5) RETURNING message_id"
	if err := db.DBconn.QueryRow(sqlStatement, newMessage.Timestamp, newMessage.FromUser, newMessage.ToUser, newMessage.ListingID, newMessage.Content).Scan(&newMessageID); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["message_id"] = strconv.Itoa(newMessageID)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
