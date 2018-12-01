package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/vnev/recyclr-backend/config"
)

func TestConnectToDB(t *testing.T) {
	var err error
	c := config.LoadConfiguration("config.json")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.DBHost, 5432, c.DBUser, c.DBPass, c.DBName)
	DBconn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Error(err.Error())
	}

	return
}
