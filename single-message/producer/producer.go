package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
)

const (
	brokers = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
)

var (
	topic = "topic"
	key   = "key"
	value = "value"
)

func main() {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "all",
		//"go.delivery.reports": "false",
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
		if ev.TopicPartition.Error != nil {
			slog.Error("Failed to deliver message: ", slog.String("error", ev.TopicPartition.Error.Error()))
			return
		}
		slog.Info("Message produced to topic", slog.String("topic", *ev.TopicPartition.Topic), slog.Int64("offset", int64(ev.TopicPartition.Offset)))
	case kafka.Error:
		slog.Error("Failed to produce message: ", slog.String("error", ev.Error()))
	}

}
