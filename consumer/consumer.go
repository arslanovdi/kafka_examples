package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"time"
)

// автономный (одиночный) потребитель (consumer)
// группа потребителей (consumer group)
// подписка на топики

// фиксация смещения асинхронная и синхронная
// автоматическая фиксация смещения... Ну это чисто конфигурация.

// десериализация сообщений

const (
	brokers    = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	asynctopic = "async-topic"
	syncTopic  = "sync-topic"
)

const (
	polltimeout = 100 * time.Millisecond
)

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "myGroup", // как минимум указать эти 2 параметра

		"auto.offset.reset": "earliest",
	})

	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println(c.GetRebalanceProtocol())

	err = c.Subscribe(asynctopic, nil)
	//err = c.SubscribeTopics([]string{asynctopic, syncTopic}, nil)

	if err != nil {
		slog.Error("Failed to subscribe to topic: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	run := true
	for run {
		msg, err := c.ReadMessage(1 * time.Second) // обертка над Poll
		if err == nil {
			fmt.Println("Message on", msg.TopicPartition, string(msg.Key), string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			slog.Error("Consumer error: ", slog.String("error", err.Error()))
		}
	}

	c.Close()
}
