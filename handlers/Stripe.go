package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	fmt.Printf("Token is %v\n", stripeToken)

	params := &stripe.ChargeParams{
		Amount:              stripe.Int64(999),
		Currency:            stripe.String(string(stripe.CurrencyUSD)),
		Description:         stripe.String("Example Charge"),
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
