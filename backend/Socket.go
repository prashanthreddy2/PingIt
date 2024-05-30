package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageStruct struct {
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Reciever  string `json:"reciever"`
	TimeStamp string `json:"timestamp"`
}

var clients = make(map[string]*websocket.Conn, 10000)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	clientId := r.URL.Query().Get("clientID")
	logrus.Infof("Websocket connected to client: [%s]", clientId)
	clients[clientId] = conn
	for {
		// Read message from WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: [%s]", err)
			return
		}
		messageQueue <- string(message)
	}
}

func BroadcastMessage() {
	for {
		select {
		case message := <-messageQueue:
			var messageS MessageStruct
			err := json.Unmarshal([]byte(message), &messageS)
			if err != nil {
				logger.Errorf("Error unmarshlling: [%s]", err)
			} else {
				logger.Infof("Message Struct: [%+v]", messageS)
				go sendMessageToClient(messageS)
				go AddMessageToDb(messageS)
			}
		}
	}
}

func sendMessageToClient(message MessageStruct) {
	senderMessage := struct {
		Message   string
		Sender    string
		Timestamp string
		Reciever  string
	}{message.Message, message.Sender, message.TimeStamp, message.Reciever}
	recieverConn, ok := clients[message.Reciever]
	if ok {
		err := recieverConn.WriteJSON(senderMessage)
		if err != nil {
			logger.Errorf("Error sending message in web socket: [%s]", err)
			return
		}
	} else {
		logger.Errorf("Reciever is not active")
		return
	}
}

func AddMessageToDb(message MessageStruct) {
	_, err := DbClient.Exec(DataEntryQuery, message.Sender, message.Reciever, message.Message, time.Now())
	if err != nil {
		logger.Errorf("Error adding message into db: [%s]", err)
		return
	} else {
		logger.Debugf("Added message to db")
	}

	_, err = DbClient.Exec(`insert into relationdb(sender, receiver) values($1, $2)`, message.Sender, message.Reciever)
	if err != nil {
		logger.Errorf("Error adding data into relationdb: [%s]", err)
	} else {
		logger.Debugf("Added relations")
	}
}
