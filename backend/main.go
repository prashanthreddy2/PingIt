package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	logger       *logrus.Logger
	DbClient     *sql.DB
	messageQueue = make(chan string)
)

func main() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	// go createProducer()
	// createConsumer()
	// time.Sleep(10 * time.Minute)
	connectToDb()
	go BroadcastMessage()
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/createAccount", handleAccountCreation)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/userList", handleUserList)
	http.HandleFunc("/fetchMessages", handleFetchMessages)
	http.HandleFunc("/fetchRelations", handleFetchRelations)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
