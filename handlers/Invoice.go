package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

//Invoice struct to hold information pertaining to an invoice
type Invoice struct {
	ID         int     `json:"invoice_id"`
	Status     bool    `json:"invoice_status"`
	Price      float64 `json:"price"`
	CreatedAt  string  `json:"created_at"`
	ForListing Listing
}

//CreateInvoice creates a new invoice and stores it into the database
//requires: Listing ID passed into request body
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var listing Listing
	_ = json.NewDecoder(r.Body).Decode(&listing)

	// confirm if the listing actually exists
	title := ""
	sqlStatement := "SELECT title FROM listings WHERE listing_id=$1"
	err := db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&title)
	if err != nil || title == "" {
		fmt.Println("CreateInvoice Query 1 fail")
		http.Error(w, "No listing found with ID", http.StatusBadRequest)
		return
	}

	sqlStatement = "INSERT INTO invoices (for_listing) VALUES ($1, $2) RETURNING invoice_id"
	invoiceID := -1
	err = db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&invoiceID)
	if err != nil || invoiceID < 0 {
		fmt.Println("CreateInvoice Query 2 fail")
		http.Error(w, "Unable to create new invoice", http.StatusInternalServerError)
		return
	}

	listingUserID, listingWeight := 0, 0
	sqlStatement = "SELECT user_id, material_weight FROM Listings WHERE id=$1"
	err = db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&listingUserID, &listingWeight)
	if err != nil {
		fmt.Println("CreateInvoice Query 3 fail")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// update points for user
	points := 50 + ((listingWeight * 10) / 100)
	sqlStatement = "UPDATE Users SET points=points+$1 WHERE user_id=$2 RETURNING points"
	err = db.DBconn.QueryRow(sqlStatement, points, listingUserID).Scan(&points)
	if err != nil {
		fmt.Println("CreateInvoice Query 4 fail")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["invoice_id"] = strconv.Itoa(invoiceID)
	resMap["new_points"] = strconv.Itoa(points)
	res, err := json.Marshal(resMap)
	if err != nil {
		fmt.Println("JSON marshal fail")
		http.Error(w, "Unable to create JSON map", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//GetInvoices returns the status and listing ID associated with
//the invoice identified by invoice_id (passed into request body)
func GetInvoices(w http.ResponseWriter, r *http.Request) {
	var invoices []Invoice
	// _ = json.NewDecoder(r.Body).Decode(&invoice)
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := `SELECT i.status, i.invoice_id, l.listing_id FROM invoices i 
					INNER JOIN Users u ON u.user_id=$1 
					INNER JOIN Listings l ON l.listing_id=i.for_listing
					ORDER BY i.created_at DESC`
	rows, err := db.DBconn.Query(sqlStatement, userID)
	if err != nil {
		http.Error(w, "Unable to query DB", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var invoice Invoice
		err = rows.Scan(&invoice.Status, &invoice.ID, &invoice.ForListing)
		invoices = append(invoices, invoice)
	}
	// Geo wants price, date, and company added below
	// resMap := make(map[string]string)
	// resMap["message"] = "Success"
	// resMap["for_listing"] = strconv.Itoa(forListingID)
	// resMap["invoice_id"] = strconv.Itoa(invoiceID)
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// res, err := json.Marshal(resMap)
	// if err != nil {
	// 	http.Error(w, "Unable to create JSON map", http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)
}

//GetInvoices returns the status and listing ID associated with
//the invoice identified by invoice_id (passed into request body)
// func GetInvoices(w http.ResponseWriter, r *http.Request) {
/* 	var invoice Invoice
   	// _ = json.NewDecoder(r.Body).Decode(&invoice)
   	params := mux.Vars(r)
   	invoiceID, err := strconv.Atoi(params["invoice_id"])
   	if err != nil {
   		http.Error(w, err.Error(), http.StatusInternalServerError)
   		return
   	}

   	forListingID := -1
   	sqlStatement := "SELECT status, for_listing FROM invoices WHERE invoice_id=$1"
   	err = db.DBconn.QueryRow(sqlStatement, invoiceID).Scan(&invoice.Status, &forListingID)
   	if err != nil {
   		http.Error(w, "Unable to query DB", http.StatusInternalServerError)
   		return
   	}

   	resMap := make(map[string]string)
   	resMap["message"] = "Success"
   	resMap["invoice_status"] = strconv.FormatBool(invoice.Status)
   	resMap["for_listing"] = strconv.Itoa(forListingID)
   	resMap["invoice_id"] = strconv.Itoa(invoiceID)

   	res, err := json.Marshal(resMap)
   	if err != nil {
   		http.Error(w, "Unable to create JSON map", http.StatusInternalServerError)
   		return
   	}
   	w.WriteHeader(http.StatusOK)
   	w.Write(res) */
// }
