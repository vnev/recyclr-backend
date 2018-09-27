package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vnev/recyclr-backend/db"
)

// GetCompanies : function to return all companies from the database
func GetCompanies(w http.ResponseWriter, r *http.Request) {
	var companies []User
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DBconn.Query("SELECT user_id, user_name FROM users where is_company=true")
	if err != nil {
		//fmt.Println(err)
		panic(err)
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
		panic(err)
	}

	json.NewEncoder(w).Encode(companies)
}

// CreateCompany : function to create a new company in the database
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var company User
	_ = json.NewDecoder(r.Body).Decode(&company)
	//fmt.Printf("read from r: addres is %s, email is %s, name is %s, pass is %s", user.Address, user.Email, user.Name, user.Password)
	sqlStatement := `
	INSERT INTO users (address, email, user_name, is_company, password)
	VALUES ($1, $2, $3, $4, crypt($5, gen_salt('md5')))
	RETURNING user_id`
	id := 0
	err := db.DBconn.QueryRow(sqlStatement, company.Address, company.Email, company.Name, true, company.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("New company created with ID:", id)
	json.NewEncoder(w).Encode(company)
}

// UpdateCompany : function to update a company in the database
func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	/*var users []User // TODO: actually get this to read in users from the DB
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var newUser User
			_ = json.NewDecoder(r.Body).Decode(&newUser)
			newUser.ID = params["id"]
			users = append(users, newUser)
			json.NewEncoder(w).Encode(newUser)
			return
		}
	}
	json.NewEncoder(w).Encode(users)*/
}

// DeleteCompany : function to delete a company from the database
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
