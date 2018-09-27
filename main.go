package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
	h "github.com/vnev/recyclr-backend/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users", h.GetUsers).Methods("GET")
	router.HandleFunc("/user", h.CreateUser).Methods("POST")

	db.ConnectToDB()
	defer db.DBconn.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}
