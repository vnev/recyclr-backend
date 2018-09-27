package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // used but unreferenced
	"github.com/vnev/recyclr-backend/config"
)

// DBconn : the main db connection object
var DBconn *sql.DB

// ConnectToDB : function that opens a connection to the database, given a config file name
func ConnectToDB() {
	var err error
	c := config.LoadConfiguration("config.json")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, 5432, c.DBUser, c.DBPass, c.DBUser)
	DBconn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = DBconn.Ping(); err != nil {
		fmt.Printf("err is %s\n", err)
		fmt.Println("Retry database connection in 5 seconds... ")
		time.Sleep(time.Duration(5) * time.Second)
		ConnectToDB()
	}

	//defer DBconn.Close()

	/*sqlStatement := `
		INSERT INTO users (age, email, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, 30, "test@test.io", "Test", "Tester").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)*/

	/*sqlStatement := "select * from users"
	var (
		id        int
		age       int
		firstName string
		lastName  string
		email     string
	)
	err = db.QueryRow(sqlStatement).Scan(&id, &age, &firstName, &lastName, &email)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Returned row is id: %d, age: %d, firstName: %s, lastName: %s, email: %s\n", id, age, firstName, lastName, email)*/

	fmt.Println("Successfully connected!")
	//return db
}
