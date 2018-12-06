package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
	h "github.com/vnev/recyclr-backend/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

//SendBackToken sends back a token for loader.io to verify that we own
// recyclr.xyz so we can run load tests
func SendBackToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "loaderio-d4781fa6082004ba4e8a3edc3dbc7299")
}

func main() {
	// The main router that handles all of our http routes
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	router.Use(loggingMiddleware)

	apiRouter.HandleFunc("/signin", h.AuthenticateUser).Methods("POST")
	apiRouter.HandleFunc("/charge", h.StripePayment).Methods("POST")
	apiRouter.HandleFunc("/user/{id}/delete", h.DeleteUser).Methods("GET")
	apiRouter.HandleFunc("/user/{id}/ban", h.BanUser).Methods("GET")

	apiRouter.HandleFunc("/user", h.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/user/rating", h.AuthMiddleware(h.UpdateRating)).Methods("PUT")
	apiRouter.HandleFunc("/user/{id}", h.AuthMiddleware(h.UpdateUser)).Methods("PUT")
	apiRouter.HandleFunc("/user/progress/{id}", h.AuthMiddleware(h.GetProgress)).Methods("GET")
	apiRouter.HandleFunc("/user/{id}", h.AuthMiddleware(h.GetUser)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/user/logout", h.AuthMiddleware(h.LogoutUser)).Methods("POST")
	apiRouter.HandleFunc("/user/transactions/{id}", h.AuthMiddleware(h.GetTransactions)).Methods("GET")
	apiRouter.HandleFunc("/user/deduct/{listing_id}", h.AuthMiddleware(h.DeductUserPoints)).Methods("POST")

	apiRouter.HandleFunc("/company", h.CreateCompany).Methods("POST")
	apiRouter.HandleFunc("/companies", h.AuthMiddleware(h.GetCompanies)).Methods("GET")
	apiRouter.HandleFunc("/company/{id}", h.AuthMiddleware(h.GetUser)).Methods("GET")                      // yes it is supposed to be GetUser not GetCompany
	apiRouter.HandleFunc("/company/transactions/{id}", h.AuthMiddleware(h.GetTransactions)).Methods("GET") // also uses the same route as users
	apiRouter.HandleFunc("/company/rating", h.AuthMiddleware(h.UpdateRating)).Methods("PUT")
	// router.HandleFunc("/company/logout", h.AuthMiddleware(h.LogoutUser)).Methods("POST")

	apiRouter.HandleFunc("/listings", h.AuthMiddleware(h.GetListings)).Methods("GET")
	apiRouter.HandleFunc("/listing/delete/{id}", h.AuthMiddleware(h.DeleteListing)).Methods("GET")
	apiRouter.HandleFunc("/listing/freeze/{id}", h.AuthMiddleware(h.FreezeListing)).Methods("POST")
	apiRouter.HandleFunc("/listing/unfreeze/{id}", h.AuthMiddleware(h.UnfreezeListing)).Methods("GET")
	apiRouter.HandleFunc("/listing/frozen/{user_id}", h.AuthMiddleware(h.GetFrozenListings)).Methods("POST")
	apiRouter.HandleFunc("/listing/{id}", h.AuthMiddleware(h.GetListing)).Methods("GET")
	apiRouter.HandleFunc("/listing", h.AuthMiddleware(h.CreateListing)).Methods("POST")
	apiRouter.HandleFunc("/listing/{id}/update", h.AuthMiddleware(h.UpdateListing)).Methods("POST")

	apiRouter.HandleFunc("/timeslots/{id}", h.AuthMiddleware(h.GetTimeslots)).Methods("GET")
	apiRouter.HandleFunc("/timeslot", h.AuthMiddleware(h.CreateTimeslot)).Methods("POST")

	apiRouter.HandleFunc("/invoice/create", h.AuthMiddleware(h.CreateInvoice)).Methods("POST")
	apiRouter.HandleFunc("/invoice/rating", h.AuthMiddleware(h.UpdateInvoiceRating)).Methods("PUT")
	apiRouter.HandleFunc("/invoice/{user_id}", h.AuthMiddleware(h.GetInvoices)).Methods("GET")

	apiRouter.HandleFunc("/messages/get", h.AuthMiddleware(h.GetMessages)).Methods("POST")
	apiRouter.HandleFunc("/messages/new", h.AuthMiddleware(h.PutMessage)).Methods("POST")

	apiRouter.HandleFunc("/loaderio-d4781fa6082004ba4e8a3edc3dbc7299.txt", SendBackToken).Methods("GET")

	db.ConnectToDB("config.json")
	defer db.DBconn.Close()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3333"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(router)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
