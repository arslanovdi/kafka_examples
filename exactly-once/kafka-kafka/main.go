package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka_examples/exactly-once/kafka-kafka/consumer"
	"kafka_examples/exactly-once/kafka-kafka/generator"
	"kafka_examples/exactly-once/kafka-kafka/model"
	"kafka_examples/exactly-once/kafka-kafka/process"
	"kafka_examples/exactly-once/kafka-kafka/producer"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	brokers = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	ticker  = 5 * time.Second
)

var (
	topic    = "topic"
	key      = "key"
	newTopic = "new-topic"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Producer
	p, err := producer.New(brokers)
	if err != nil {
		slog.Error("Failed to create producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer p.Close()

	// Producer with transactional id
	ptrans, err := producer.New(
		brokers,
		producer.WithTransaction("my-transactional-id"),
	)
	if err != nil {
		slog.Error("Failed to create transactional producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer ptrans.Close()

	// consumer withOUT auto commit
	c, err := consumer.New(brokers, "group")
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer c.Close()

	// consumer with auto commit
	cauto, err := consumer.New(brokers, "group-new-msg", consumer.WithAutoCommit())
	if err != nil {
		slog.Error("Failed to create autocommit consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cauto.Close()

	go generator.Generate(p, topic, key, ticker) // генерация сообщений со случайными числами

	go func() { // чтение, обработка и запись
		err = c.Run(topic, ptrans, process.Multiplier)
		if err != nil {
			slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	go func() { // чтение новых сообщений
		err = cauto.Run(newTopic, nil, func(p model.Producer, msg *kafka.Message, cmsg producer.CommitMsgData) {
			fmt.Printf("Read New Message from topic: %s. Value: %s\n", newTopic, msg.Value)
			//fmt.Printf("message\nkey:%s\nvalue:%s\noffset:%d\n\n", msg.Key, msg.Value, msg.TopicPartition.Offset)
		})
		if err != nil {
			slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-stop
}
