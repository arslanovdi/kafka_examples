package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka_examples/rebalance/consumer"
	"kafka_examples/rebalance/generator"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	brokers        = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	group          = "Mygroup"
	topic          = "test"
	key            = "key"
	partitionCount = 10
	consumerCount  = 7
)

func main() {

	// Создание топика с 10-ю партициями
	aClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		slog.Error("Failed to create admin client: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer aClient.Close()
	spec := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     partitionCount,
		ReplicationFactor: 3,
	}
	_, err = aClient.CreateTopics(context.Background(), []kafka.TopicSpecification{spec})
	if err != nil {
		slog.Error("Failed to create topic: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Генерация сообщений
	for i := range partitionCount {
		go generator.Generate(brokers, topic, key, 5*time.Second, int32(i))
	}

	// запускаем консюмеры
	consumers := make([]*consumer.Rebalanceconsumer, consumerCount)

	for i := range consumerCount {
		consumers[i] = consumer.New(brokers, group)
		defer consumers[i].Close()
		go func() {
			err := consumers[i].Run(topic)
			if err != nil {
				slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
				os.Exit(1)
			}
		}()
	}

	// через 20 сек отключаем консюмер и смотрим в логах как происходит ребалансировка
	time.Sleep(20 * time.Second)

	consumers[consumerCount/2].Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}
