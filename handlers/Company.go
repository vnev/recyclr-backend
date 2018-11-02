package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vnev/recyclr-backend/db"
)

// GetCompanies returns all companies from the database in the JSON format. It does not require
// any parameters.
func GetCompanies(w http.ResponseWriter, r *http.Request) {
	var companies []User
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DBconn.Query("SELECT user_id, user_name FROM users where is_company=true")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var company User
		err = rows.Scan(&company.ID, &company.Name)
		fmt.Printf("ID is %d, Name is %s\n", company.ID, company.Name)
		companies = append(companies, company)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(companies)
}

// CreateCompany creates a new company in the database, and returns the newly created company in JSON format.
// In the request body, it expects an address, email, user_name, is_company, and password.
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var company User
	_ = json.NewDecoder(r.Body).Decode(&company)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO users (address, city, state, email, user_name, is_company, passwd)
	VALUES ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('md5')))
	RETURNING user_id`
	id := 0
	err := db.DBconn.QueryRow(sqlStatement, company.Address, company.City, company.State, company.Email, company.Name, true, company.Password).Scan(&id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("New company created with ID:", id)
	json.NewEncoder(w).Encode(company)
}

// UpdateCompany updates a company in the database, and will return the number of rows updated. It expects
// the user_id, along with all the fields that are requesting to be changed with their new information.
func UpdateCompany(w http.ResponseWriter, r *http.Request) {

}

// DeleteCompany deletes a company from the database. It expects the user_id, and will only work if
// the user sending the request has sufficient admin priveliges.
func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	/*var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)*/
}
