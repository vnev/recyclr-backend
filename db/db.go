// Package db implements the connection to our PostgreSQL server.
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // used but unreferenced
	"github.com/vnev/recyclr-backend/config"
)

// DBconn is the main database connection object, used globally.
var DBconn *sql.DB

// ConnectToDB opens a connection to the database, and keeps it open while the server is running.
func ConnectToDB() {
	var err error
	c := config.LoadConfiguration("config.json")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.DBHost, 5432, c.DBUser, c.DBPass, c.DBName)
	DBconn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = DBconn.Ping(); err != nil {
		fmt.Printf("err is %s\n", err)
		fmt.Println("Retrying database connection in 5 seconds...")
		time.Sleep(time.Duration(5) * time.Second)
		ConnectToDB()
	}

	fmt.Println("Successfully connected to database!")
}
