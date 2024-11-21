package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	brokers               = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	topic                 = "topic"
	group                 = "MyGroup"
	defaultSessionTimeout = 6000
	ReadTimeout           = 1 * time.Second
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       brokers,
		"group.id":                group,
		"session.timeout.ms":      defaultSessionTimeout,
		"enable.auto.commit":      "true",
		"auto.commit.interval.ms": 1000,
		"auto.offset.reset":       "earliest", // earliest - сообщения с commit offset; latest - новые сообщение
	})
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		slog.Error("Failed to subscribe to topic: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-stop:
			break loop
		default:
			msg, err := consumer.ReadMessage(ReadTimeout) // read message with timeout
			if err == nil {
				fmt.Printf("Message on topic %s: key = %s, value = %s\n", *msg.TopicPartition.Topic, string(msg.Key), string(msg.Value))
			} else if !err.(kafka.Error).IsTimeout() { // TODO process timeout
				slog.Error("Consumer error: ", slog.String("error", err.Error()))
			}
		}
	}
}
