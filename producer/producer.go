package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"time"
)

// отправка сообщений всегда асинхронная, есть отчеты о доставке
// сериализация сообщений, пользовательский сериализатор написать.
// работа с Headers
// интерцепторы, написание своего?

const (
	brokers    = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	asynctopic = "async-topic"
	syncTopic  = "sync-topic"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{ // https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
		"bootstrap.servers": brokers, // список адресов брокеров host или host:port; должен быть указан

		"client.id":       "myServiceId", // идентификатор клиента
		"retention.bytes": -1,
		"retention.ms":    7 * 24 * time.Hour / time.Millisecond, // default: 7 дней
		"acks":            "all",                                 // ожидать подтверждения от всех синхронизированных реплик (none, one, all). Контролировать min.insync.replicas. Если живых реплик будет < min.insync.replicas, то сообщение не будет доставлено!
		//"go.delivery.reports": "false",       // отключить отчеты о доставке
	})

	if err != nil {
		slog.Error("Failed to create producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer p.Close()

	// читаем события из канала Events
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					slog.Error("Failed to deliver message: ", slog.String("error", ev.TopicPartition.Error.Error()))
					continue
				}
				fmt.Println("Message delivered to topic", ev.TopicPartition)
			}
		}
	}()

	topic := asynctopic
	// Асинхронная отправка сообщений в топик
	for _, word := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
				/*Offset:      0,
				Metadata:    nil,
				Error:       nil,
				LeaderEpoch: nil,*/
			}, Value: []byte(word)},
			nil)
	}

	topic = syncTopic

	p.Flush(15 * 1000) // блокируемся пока не будут доставлены сообщения либо не истечет тайм-аут
}
