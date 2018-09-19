package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/vnev/recyclr-backend/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user", h.GetUser).Methods("GET")
	router.HandleFunc("/user", h.CreateUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
