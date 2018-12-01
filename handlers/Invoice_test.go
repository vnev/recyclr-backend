package handlers

import (
	"strconv"
	"testing"

	"github.com/vnev/recyclr-backend/db"
)

//CreateInvoice creates a new invoice and stores it into the database
//requires: Listing ID passed into request body
/* func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("body is %v", r.Body)
	var listing Listing
	_ = json.NewDecoder(r.Body).Decode(&listing)

	fmt.Printf("listingID passed in is %d\n", listing.ID)

	// confirm if the listing actually exists
	title := ""
	sqlStatement := "SELECT title FROM listings WHERE listing_id=$1"
	err := db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&title)
	if err != nil || title == "" {
		fmt.Println("CreateInvoice Query 1 fail")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement = "INSERT INTO invoices (for_listing) VALUES ($1) RETURNING invoice_id"
	invoiceID := -1
	err = db.DBconn.QueryRow(sqlStatement, listing.ID).Scan(&invoiceID)
	if err != nil || invoiceID < 0 {
		fmt.Println("CreateInvoice Query 2 fail")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listingUserID, listingWeight := 0, 0
	sqlStatement = "SELECT user_id, material_weight FROM Listings WHERE listing_id=$1"
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
} */

//GetInvoices returns the status and listing ID associated with
//the invoice identified by invoice_id (passed into request body)
func TestGetInvoices(t *testing.T) {
	type subinvoice struct {
		ID          int     `json:"invoice_id"`
		Status      bool    `json:"invoice_status"`
		Price       float64 `json:"price"`
		CreatedAt   string  `json:"created_at"`
		CompanyName string  `json:"company_name"`
		UserName    string  `json:"user_name"`
		Title       string  `json:"title"`
	}

	var invoices []subinvoice
	//params := mux.Vars(r)
	userID, err := strconv.Atoi("1")
	if err != nil {
		t.Error(err.Error())
	}

	sqlStatement := `SELECT i.status, i.invoice_id, l.price, l.title, u.user_name, u2.user_name, i.created_at
					FROM invoices i 
					INNER JOIN Users u ON u.user_id=$1 
					INNER JOIN Listings l ON l.listing_id=i.for_listing
					INNER JOIN Users u2 ON l.frozen_by=u2.user_id
					WHERE l.user_id=$2
					ORDER BY i.created_at DESC`
	rows, err := db.DBconn.Query(sqlStatement, userID, userID)
	if err != nil {
		t.Error(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var invoice subinvoice
		err = rows.Scan(&invoice.Status, &invoice.ID, &invoice.Price, &invoice.Title, &invoice.UserName, &invoice.CompanyName, &invoice.CreatedAt)
		invoices = append(invoices, invoice)
	}
	if err = rows.Err(); err != nil {
		t.Error(err.Error())
	}

}
