package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/vnev/recyclr-backend/db"
)

// User struct contains the user schema in a struct format.
type User struct {
	ID        int    `json:"user_id"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IsCompany bool   `json:"is_company"`
	Rating    int    `json:"rating"`
	JoinedOn  string `json:"joined_on"`
	Password  string `json:"passwd"`
	Token     string `json:"token"`
	City      string `json:"city"`
	State     string `json:"state"`
	Points    int    `json:"points"`
}

// GetUser returns a user from the database in JSON format, given the specific user_id as a URL parameter.
// TODO: maybe just return a newly defined struct without password field.
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := "SELECT user_id, points, address, city, state, email, user_name, is_company, rating, joined_on FROM users WHERE user_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Points, &user.Address, &user.City, &user.State, &user.Email, &user.Name, &user.IsCompany, &user.Rating, &user.JoinedOn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&user)
}

// CreateUser creates a new user in the database. It expects address, email, user_name, is_company, and passwd fields.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO users (address, city, state, email, user_name, is_company, passwd)
	VALUES ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('md5')))
	RETURNING user_id`
	id := 0
	err := db.DBconn.QueryRow(sqlStatement, user.Address, user.City, user.State, user.Email, user.Name, false, user.Password).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("New user created with ID:", id)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user in the database, given its' user_id and other fields requesting to be changed.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	/*if user.ID == 0 {
		fmt.Println("Bad request 1")
		http.Error(w, "No user ID found", http.StatusBadRequest)
		return
	}*/

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		bearer := strings.Split(authHeader, " ")
		if len(bearer) == 2 {
			token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error")
				}
				return []byte("secret"), nil
			})
			if err != nil {
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			if token.Valid {
				// do nothing I guess?
				fmt.Println("should be valid, continuing...")
				token := ""
				sqlStatement := "SELECT token FROM users WHERE user_id=$1"
				err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&token)

				if bearer[1] != token {
					// trying to modify other user, reject this
					http.Error(w, "Trying to modify another user", http.StatusInternalServerError)
					return
				}
			} else {
				resMap := make(map[string]string)
				resMap["message"] = "Failed"

				res, err := json.Marshal(resMap)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res)
				return
			}
		} else {
			http.Error(w, "Invalid authorization header", http.StatusBadRequest)
			return
		}
	} else {
		resMap := make(map[string]string)
		resMap["message"] = "Failed"

		res, err := json.Marshal(resMap)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Bad request 2")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE users SET "

	// Time to iteratively loop over a struct, the easiest to understand syntax ever!
	structIterator := reflect.ValueOf(user)
	for i := 0; i < structIterator.NumField(); i++ {
		//fmt.Printf("field: %+v, value: %+v\n", structIterator.Type().Field(i).Name, structIterator.Field(i).Interface())
		field := structIterator.Type().Field(i).Name
		/*fmt.Printf("Field is %s\n", field)
		if field != "Address" {
			fmt.Printf("not address\n")
			if field != "Email" {
				fmt.Printf("not email\n")
				if field != "Name" {
					fmt.Printf("not name\n")
					if field != "Password" {
						fmt.Printf("not passwd\n")
						continue
					}
				}
			}
		}*/

		val := structIterator.Field(i).Interface()

		// Check if the field is zero-valued, meaning it won't be updated
		//fmt.Printf("VAL IS %v and TYPE IS %v and ZERO OF TYPE IS %v\n", val, structIterator.Field(i).Type(), reflect.Zero(structIterator.Field(i).Type()).Interface())
		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			// fmt.Printf("%v is non-zero, adding to update\n", field)
			if strings.ToLower(field) == "name" {
				sqlStatement = sqlStatement + "user_name=$" + strconv.Itoa(j) + ", "
			} else if strings.ToLower(field) == "password" {
				// crypt($5, gen_salt('md5'))
				sqlStatement = sqlStatement + "passwd=crypt($" + strconv.Itoa(j) + ", gen_salt('md5')), "
			} else {
				sqlStatement = sqlStatement + strings.ToLower(field) + "=$" + strconv.Itoa(j) + ", "
			}

			j++
			values = append(values, val)
		}
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-2]
	sqlStatement = sqlStatement + " WHERE user_id" + "=$" + strconv.Itoa(j)
	values = append(values, userID)
	fmt.Printf("executing SQL: \n\t%s\n", sqlStatement)
	// fmt.Printf("$1 is %s and $2 is %d\n", values[0], values[1])
	row, err := db.DBconn.Exec(sqlStatement, values...) //.Scan(&user.ID, &user.Name)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(count, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	//json.NewEncoder(w).Encode({"status": "200", "message": "success"})
}

// BanUser bans a specific user given their user_id as a URL parameter.
func BanUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := "UPDATE users SET banned='t' WHERE user_id=$1"
	_, err = db.DBconn.Exec(sqlStatement, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Successfully banned"
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// AuthenticateUser generates a JWT for the user and returns it in JSON format.
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	sqlStatement := "SELECT email FROM users WHERE email=$1 AND banned='f'"
	email := ""
	_ = db.DBconn.QueryRow(sqlStatement, user.Email).Scan(&email)
	if email == "" {
		http.Error(w, "No user found with email "+user.Email, http.StatusBadRequest)
		return
	}

	userID := 0
	userName := ""

	sqlStatement = "SELECT user_id, user_name FROM users WHERE email=$1 AND passwd=crypt($2, passwd)"
	err := db.DBconn.QueryRow(sqlStatement, user.Email, user.Password).Scan(&userID, &userName)
	if err != nil {
		fmt.Printf("Err is %s\n", err.Error())
		http.Error(w, "No user found with that email/password", http.StatusBadRequest)
		return
	}
	// if err != nil {
	// 	http.Error(w, err.Error(), HTTPInternalError)
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "recyclr.xyz",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"name": user.Name,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement = "UPDATE users SET token=$1 WHERE user_id=$2"
	_, err = db.DBconn.Exec(sqlStatement, tokenString, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["token"] = tokenString
	resMap["user_id"] = strconv.Itoa(userID)
	resMap["user_name"] = userName

	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// LogoutUser logs a user out, setting their JWT to 0. It expects the user_id to be sent in the request body.
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.ID == 0 {
		http.Error(w, "No user ID found", http.StatusBadRequest)
		return
	}

	sqlStatement := "UPDATE users SET token='0' WHERE user_id=$1"
	_, err := db.DBconn.Exec(sqlStatement, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteUser deletes a listing from the database given their user_id. It will only work if
// the user sending the request has sufficient admin priveliges.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error. Couldn't parse user ID", http.StatusInternalServerError)
		return
	}

	sqlStatement := "DELETE FROM users WHERE user_id=$1"
	_, err = db.DBconn.Exec(sqlStatement, userID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Successfully deleted user"
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetProgress gets the progress of a user's listings, returning it in JSON format given their user_id as a URL parameter.
func GetProgress(w http.ResponseWriter, r *http.Request) {
	var listings []Listing
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.DBconn.Query("SELECT listing_id, title, description, material_type, material_weight, active FROM listings WHERE user_id=$1", userID)
	//err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Address, &user.Email, &user.Name, &user.IsCompany, &user.Rating, &user.JoinedOn)

	defer rows.Close()
	for rows.Next() {
		var listing Listing
		err = rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.MaterialType, &listing.MaterialWeight, &listing.Active)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		listings = append(listings, listing)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listings)
}

// GetTransactions returns all orders for a company in JSON format, given their user_id as a URL parameter.
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.DBconn.Query("SELECT order_id, user_id, company_id, total, confirmed FROM orders WHERE user_id=$1", userID)
	//err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Address, &user.Email, &user.Name, &user.IsCompany, &user.Rating, &user.JoinedOn)

	defer rows.Close()
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.UserID, &order.CompanyID, &order.Total, &order.Confirmed)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
