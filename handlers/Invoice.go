package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

//Invoice struct to hold information pertaining to an invoice
type Invoice struct {
	ID         int  `json:"invoice_id"`
	Status     bool `json:"invoice_status"`
	ForListing Listing
}

//CreateInvoice creates a new invoice and stores it into the database
//requires: Listing ID passed into request body
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var listing Listing
	_ = json.NewDecoder(r.Body).Decode(&listing)

	title := ""
	sqlStatement := "SELECT title FROM listings WHERE listing_id=$1"
	err := db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&title)
	if err != nil || title == "" {
		http.Error(w, "No listing found with ID", http.StatusBadRequest)
		return
	}

	sqlStatement = "INSERT INTO invoices (for_listing) VALUES ($1) RETURNING invoice_id"
	invoiceID := -1
	err = db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&invoiceID)
	if err != nil || invoiceID < 0 {
		http.Error(w, "Unable to create new invoice", http.StatusInternalServerError)
		return
	}

	listingUserID, listingWeight := 0, 0
	sqlStatement = "SELECT user_id, material_weight FROM Listings WHERE id=$1"
	err = db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&listingUserID, &listingWeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// update points for user
	points := 50 + ((listingWeight * 10) / 100)
	sqlStatement = "UPDATE Users SET points=points+$1 WHERE user_id=$2 RETURNING points"
	err = db.DBconn.QueryRow(sqlStatement, points, listingUserID).Scan(&points)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["invoice_id"] = strconv.Itoa(invoiceID)
	resMap["new_points"] = strconv.Itoa(points)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, "Unable to create JSON map", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//GetInvoice returns the status and listing ID associated with
//the invoice identified by invoice_id (passed into request body)
func GetInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice Invoice
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
	w.Write(res)
}
