package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func createProducer() {
	message := fmt.Sprintf("Hello from producer at time: [%s]", time.Now())
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka-headless.kafka:9092"},
	})
	ctx := context.Background()
	for {
		err := writer.WriteMessages(ctx, kafka.Message{
			Topic: "test",
			Value: []byte(message),
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
		Brokers:   []string{"kafka-headless.kafka:9092"},
		Topic:     "test",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		GroupID:     "test-consumer-group",
		StartOffset: kafka.LastOffset,
		CommitInterval: time.Second,
	})
	ctx := context.Background()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			break
		}
		logger.Infof("Message at offset %d: [%s]", m.Offset, string(m.Value))
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
