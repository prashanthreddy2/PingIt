package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
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
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-headless.kafka:9092"},
	})
	ctx := context.Background()
	for {
		// Read message from WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: [%s]", err)
			return
		}
		err = writer.WriteMessages(ctx, kafka.Message{
			Topic: "test",
			Value: message,
		})
		if err != nil {
			logger.Errorf("Error writing from kafka producer: [%s]", err)
			return
		}
	}
}

func BroadcastMessage() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka-headless.kafka:9092"},
		Topic:          "test",
		Partition:      0,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		GroupID:        "test-consumer-group",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
	ctx := context.Background()
	var message MessageStruct
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			break
		}
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			logger.Errorf("Error unmarshlling the message: [%s]", err)
			logger.Infof("Message at offset %d: [%s]", m.Offset, string(m.Value))
		} else {
			logger.Infof("Message at offset %d: [%+v]", m.Offset, message)
		}
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			logger.Errorf("Error unmarshlling the message: [%s]", err)
			logger.Infof("Message at offset %d: [%s]", m.Offset, string(m.Value))
		} else {
			logger.Infof("Message at offset %d: [%+v]", m.Offset, message)
		}
		go sendMessageToClient(message)
		go AddMessageToDb(message)
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
