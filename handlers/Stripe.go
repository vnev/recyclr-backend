package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vnev/recyclr-backend/db"

	"github.com/stripe/stripe-go/charge"
	"github.com/vnev/recyclr-backend/config"

	"github.com/stripe/stripe-go"
)

//StripePayment handles a payment from Stripe, given the payment token in a form in the request body.
func StripePayment(w http.ResponseWriter, r *http.Request) {
	c := config.LoadConfiguration("config.json")
	stripe.Key = c.StripeSecret
	//stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
	stripeToken := r.FormValue("token")
	listingID, err := strconv.Atoi(r.FormValue("listing_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type attributes struct {
		Price float64
		Title string
		Email string
	}

	var attr attributes
	sqlStatement := "SELECT l.price, l.title, u.email FROM Listings l INNER JOIN Users u ON l.user_id=u.user_id WHERE l.listing_id=$1"
	if err := db.DBconn.QueryRow(sqlStatement, listingID).Scan(&attr.Price, &attr.Title, &attr.Email); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Token is %v\n", stripeToken)

	params := &stripe.ChargeParams{
		Amount:              stripe.Int64(int64(attr.Price * 100)),
		Currency:            stripe.String(string(stripe.CurrencyUSD)),
		Description:         stripe.String("Listing " + attr.Title + " by " + attr.Email),
		StatementDescriptor: stripe.String("Recyclr"),
	}
	params.SetSource(stripeToken)
	ch, err := charge.New(params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["stripe_status"] = ch.Status

	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
