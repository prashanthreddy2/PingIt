package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "10.99.144.3"
	port     = 5432
	user     = "ping-it-backend	"
	password = "Prashanth@123"
	dbname   = "postgres"
)

func connectToDb() {
	if DbClient != nil {
		return
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Errorf("Error creating db client: [%s]", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		logger.Errorf("Error pinging db: [%s]", err)
		return
	}
	// defer db.Close()
	DbClient = db
	logger.Infof("Connected to db")
}
