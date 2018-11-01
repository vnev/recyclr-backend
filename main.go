package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
	h "github.com/vnev/recyclr-backend/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		// w.Header().Set("Access-Control-Allow-Origin", "*")
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

	/* router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
	}) */

	router.Use(loggingMiddleware)

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
	router.HandleFunc("/user/deduct/{listing_id}", h.AuthMiddleware(h.DeductUserPoints)).Methods("POST")

	router.HandleFunc("/company", h.CreateCompany).Methods("POST")
	router.HandleFunc("/companies", h.AuthMiddleware(h.GetCompanies)).Methods("GET")
	router.HandleFunc("/company/{id}", h.AuthMiddleware(h.GetUser)).Methods("GET") // yes it is supposed to be GetUser not GetCompany
	router.HandleFunc("/company/delete", h.AuthMiddleware(h.DeleteCompany)).Methods("POST")
	router.HandleFunc("/company/transactions/{id}", h.AuthMiddleware(h.GetTransactions)).Methods("GET") // also uses the same route as users
	// router.HandleFunc("/company/logout", h.AuthMiddleware(h.LogoutUser)).Methods("POST")

	router.HandleFunc("/listings", h.AuthMiddleware(h.GetListings)).Methods("GET")
	router.HandleFunc("/listing/freeze/{id}", h.AuthMiddleware(h.FreezeListing)).Methods("POST")
	router.HandleFunc("/listing/unfreeze/{id}", h.AuthMiddleware(h.UnfreezeListing)).Methods("GET")
	router.HandleFunc("/listing/frozen/{user_id}", h.AuthMiddleware(h.GetFrozenListings)).Methods("POST")
	router.HandleFunc("/listing/{id}", h.AuthMiddleware(h.GetListing)).Methods("GET")
	router.HandleFunc("/listing", h.AuthMiddleware(h.CreateListing)).Methods("POST")
	router.HandleFunc("/listing/{id}/update", h.AuthMiddleware(h.UpdateListing)).Methods("POST")

	router.HandleFunc("/timeslots/{id}", h.AuthMiddleware(h.GetTimeslots)).Methods("GET")
	router.HandleFunc("/timeslot", h.AuthMiddleware(h.CreateTimeslot)).Methods("POST")

	router.HandleFunc("/invoice/create", h.AuthMiddleware(h.CreateInvoice)).Methods("POST")
	router.HandleFunc("/invoice/{user_id}", h.AuthMiddleware(h.GetInvoices)).Methods("GET")

	router.HandleFunc("/messages/get", h.AuthMiddleware(h.GetMessages)).Methods("POST")
	router.HandleFunc("/messages/new", h.AuthMiddleware(h.PutMessage)).Methods("POST")

	router.HandleFunc("/loaderio-d4781fa6082004ba4e8a3edc3dbc7299.txt", SendBackToken).Methods("GET")

	db.ConnectToDB()
	defer db.DBconn.Close()

	// handler := cors.Default().Handler(router)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
