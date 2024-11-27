package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
)

const (
	brokers = "localhost:29092,localhost:29093,localhost:29094"
)

var (
	topic = "topic"
	key   = "key"
	value = "value"
)

func main() {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,

		"security.protocol":        "ssl",
		"ssl.key.location":         "security/SSL/PEM/certs/producer.key",        // Path to client's private key (PEM) used for authentication.
		"ssl.certificate.location": "security/SSL/PEM/certs/producer-signed.crt", // Path to client's public key (PEM) used for authentication.
		"ssl.ca.location":          "security/SSL/PEM/certs/rootCA.crt",          // Path to CA certificate.
		//"enable.ssl.certificate.verification": false, // don't use in production
		"acks": "all",
	})
	if err != nil {
		slog.Error("Failed to create producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            []byte(key),
		Value:          []byte(value),
	}, deliveryChan)

	if err != nil {
		slog.Error("Failed to produce message: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	event := <-deliveryChan
	switch ev := event.(type) {
	case *kafka.Message:
		slog.Info("Message produced to topic", slog.String("topic", *ev.TopicPartition.Topic), slog.Int64("offset", int64(ev.TopicPartition.Offset)))
	case kafka.Error:
		slog.Error("Failed to produce message: ", slog.String("error", ev.Error()))
	}

}
