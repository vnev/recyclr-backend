package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vnev/recyclr-backend/db"
	h "github.com/vnev/recyclr-backend/handlers"
)

func main() {
	router := mux.NewRouter()
	// router.HandleFunc("/users", h.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/user", h.CreateUser).Methods("POST")

	router.HandleFunc("/companies", h.GetCompanies).Methods("GET")
	router.HandleFunc("/company/{id}", h.GetUser).Methods("GET") // yes it is supposed to be GetUser not GetCompany
	router.HandleFunc("/company", h.CreateCompany).Methods("POST")

	router.HandleFunc("/listings", h.GetListings).Methods("GET")
	router.HandleFunc("/listing/{id}", h.GetListing).Methods("GET")
	router.HandleFunc("/listing", h.CreateListing).Methods("POST")
	router.HandleFunc("/listing/update", h.UpdateListing).Methods("POST")

	router.HandleFunc("/signin", h.AuthenticateUser).Methods("POST")

	db.ConnectToDB()
	defer db.DBconn.Close()

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
