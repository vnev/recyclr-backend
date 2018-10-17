package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vnev/recyclr-backend/db"
	h "github.com/vnev/recyclr-backend/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/signin", h.AuthenticateUser).Methods("POST")
	router.HandleFunc("/charge", h.StripePayment).Methods("POST")
	router.HandleFunc("/user/{id}/delete", h.DeleteUser).Methods("GET")
	router.HandleFunc("/user/{id}/ban", h.BanUser).Methods("GET")

	router.HandleFunc("/user", h.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", h.AuthMiddleware(h.UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/progress/{id}", h.AuthMiddleware(h.GetProgress)).Methods("GET")
	router.HandleFunc("/user/{id}", h.AuthMiddleware(h.GetUser)).Methods("GET")
	router.HandleFunc("/user/logout", h.AuthMiddleware(h.LogoutUser)).Methods("POST")
	router.HandleFunc("/user/transactions/{id}", h.AuthMiddleware(h.GetTransactions)).Methods("GET")

	router.HandleFunc("/company", h.CreateCompany).Methods("POST")
	router.HandleFunc("/companies", h.AuthMiddleware(h.GetCompanies)).Methods("GET")
	router.HandleFunc("/company/{id}", h.AuthMiddleware(h.GetUser)).Methods("GET") // yes it is supposed to be GetUser not GetCompany
	router.HandleFunc("/company/delete", h.AuthMiddleware(h.DeleteCompany)).Methods("POST")
	router.HandleFunc("/company/transactions/{id}", h.AuthMiddleware(h.GetTransactions)).Methods("GET") // also uses the same route as users
	// router.HandleFunc("/company/logout", h.AuthMiddleware(h.LogoutUser)).Methods("POST")

	router.HandleFunc("/listings", h.AuthMiddleware(h.GetListings)).Methods("GET")
	router.HandleFunc("/listing/{id}", h.AuthMiddleware(h.GetListing)).Methods("GET")
	router.HandleFunc("/listing", h.AuthMiddleware(h.CreateListing)).Methods("POST")
	router.HandleFunc("/listing/{id}/update", h.AuthMiddleware(h.UpdateListing)).Methods("POST")

	router.HandleFunc("/timeslots/{id}", h.AuthMiddleware(h.GetTimeslots)).Methods("GET")
	router.HandleFunc("/timeslot", h.AuthMiddleware(h.CreateTimeslot)).Methods("POST")

	db.ConnectToDB()
	defer db.DBconn.Close()

	handler := cors.Default().Handler(router)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
