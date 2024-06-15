package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaTestStruct struct {
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
	Reciever  string `json:"reciever"`
}

func createProducer() {
	senderMessage := kafkaTestStruct{
		Message:   "message",
		Sender:    "sender",
		Timestamp: "timestamp",
		Reciever:  "reciever",
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-headless.kafka:9092"},
	})
	ctx := context.Background()
	message, err := json.Marshal(senderMessage)
	if err != nil {
		logger.Errorf("Error marshlling the message: [%s]", err)
		message = []byte("test message")
	}
	for {
		err := writer.WriteMessages(ctx, kafka.Message{
			Topic: "test",
			Value: message,
		})
		if err != nil {
			logger.Errorf("Error creating the kafka consumer: [%s]", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	if err := writer.Close(); err != nil {
		logger.Errorf("Error closing the producer: [%s]", err)
	} else {
		logger.Infof("Wrote the message: [%s]", message)
	}

}

func createConsumer() {
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
			logger.Errorf("Error reading message: %s", err)
			break
		}
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			logger.Errorf("Error unmarshlling the message: [%s]", err)
			logger.Infof("Message at offset %d: [%s]", m.Offset, string(m.Value))
		} else {
			logger.Infof("Message at offset %d: [%+v]", m.Offset, message)
		}
		// if err := reader.CommitMessages(ctx, m); err != nil {
		// 	logger.Errorf("Failed to commit message: %v", err)
		// }
	}
	if err := reader.Close(); err != nil {
		logger.Errorf("Error closing the consumer: %v", err)
	} else {
		logger.Infof("Consumer closed")
	}
}
