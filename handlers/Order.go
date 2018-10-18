package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

//Order struct contains the order DB schema in a struct format
type Order struct {
	ID        int  `json:"order_id"`
	UserID    int  `json:"user_id"`
	CompanyID int  `json:"company_id"`
	Total     int  `json:"total"`
	Confirmed bool `json:"confirmed"`
}

// GetOrder returns an order from the database in JSON format, given the specific order_id as a URL parameter.
func GetOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	orderID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT user_id, company_id, total, confirmed FROM orders WHERE order_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, orderID).Scan(&order.UserID, &order.CompanyID, &order.Total, &order.Confirmed)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&order)
}

// GetOrders returns all orders from the database in JSON format for a specific given user,
// given their user_id as a URL parameter.
func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := db.DBconn.Query("SELECT user_id, company_id, total, confirmed FROM orders WHERE user_id=$1", userID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.UserID, &order.CompanyID, &order.Total, &order.Confirmed)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// CreateOrder creates a new listing in the database. It expects user_id, company_id, total, and confirmed.
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}
	//fmt.Println("LISTING IS: ", listing)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO orders (user_id, company_id, total, confirmed)
	VALUES ($1, $2, $3, $4)
	RETURNING order_id`
	id := 0
	err = db.DBconn.QueryRow(sqlStatement, order.UserID, order.CompanyID, order.Total, order.Confirmed).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("New order created with ID: ", id)
	json.NewEncoder(w).Encode(order)
}

// UpdateOrder updates an order in the database, given its' order_id and other fields requesting to be changed.
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request parameters", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE orders SET "

	structIterator := reflect.ValueOf(order)
	for i := 0; i < structIterator.NumField(); i++ {
		field := structIterator.Type().Field(i).Name
		val := structIterator.Field(i).Interface()

		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			sqlStatement += strings.ToLower(field) + "=$" + strconv.Itoa(j) + ", "
			j++
			values = append(values, val)
		}
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-2]
	sqlStatement = sqlStatement + " WHERE order_id" + "=$" + strconv.Itoa(j)
	values = append(values, orderID)
	row, err := db.DBconn.Exec(sqlStatement, values...)
	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteOrder deletes an order from the database given its' order_id. It will only work if
// the user sending the request has sufficient admin priveliges.
func DeleteOrder(w http.ResponseWriter, r *http.Request) {

}
