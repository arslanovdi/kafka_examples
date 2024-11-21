package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	brokers = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	tick    = 5 * time.Millisecond
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
		"go.batch.producer":  true,
		"linger.ms":          1000,
		"batch.num.messages": 10000,
		"batch.size":         1000000,
		"message.max.bytes":  1000000,

		"compression.codec": "lz4",
		"compression.level": -1,
	})
	if err != nil {
		slog.Error("Failed to create producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event, 1000)
	defer close(deliveryChan)

	stopevents := make(chan bool)
	stopproduce := make(chan bool)
	var wgProduce sync.WaitGroup

	// events
	go func() {
	loop:
		for {
			select {
			case <-stopevents:
				break loop
			case event := <-deliveryChan:
				switch ev := event.(type) {
				case *kafka.Message:
					slog.Info("Message produced to topic", slog.String("topic", *ev.TopicPartition.Topic), slog.Int64("offset", int64(ev.TopicPartition.Offset)))
				case kafka.Error:
					slog.Error("Failed to produce message: ", slog.String("error", ev.Error()))
				}
			}
		}
	}()

	// produce
	wgProduce.Add(1)
	go func() {
		defer wgProduce.Done()
		ticker := time.NewTicker(tick)
	loop:
		for {
			select {
			case <-stopproduce:
				break loop
			case <-ticker.C:
				err = producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic},
					Key:            []byte(key),
					Value:          []byte(value),
				}, deliveryChan)

				if err != nil {
					slog.Error("Failed to produce message: ", slog.String("error", err.Error()))
					os.Exit(1)
				}
			}
		}
		producer.Flush(1000)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	stopproduce <- true
	wgProduce.Wait()
	stopevents <- true
}
