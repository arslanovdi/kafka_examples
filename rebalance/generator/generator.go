package generator

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Generate добавляет в партицию № partition топика topic случайное число с шагом ticker
func Generate(brokers string, topic string, key string, ticker time.Duration, partition int32) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		slog.Error("Failed to create producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	tik := time.NewTicker(ticker)

loop:
	for {
		select {
		case <-stop:
			break loop
		case <-tik.C:
			_ = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},

				Key:   []byte(key),
				Value: []byte(fmt.Sprintf("%d", rand.Int31())),
			}, nil)
		}
	}
}
