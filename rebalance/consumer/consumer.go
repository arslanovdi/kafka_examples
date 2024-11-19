package consumer

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"sync"
	"time"
)

const (
	ReadTimeout           = 1 * time.Second
	defaultSessionTimeout = 6000
)

type Rebalanceconsumer struct {
	consumer *kafka.Consumer
	stop     bool           // command to stop reading messages
	readers  sync.WaitGroup // runned readers
}

// Run reads messages from kafka topic
func (rb *Rebalanceconsumer) Run(topic string) error {
	if rb.stop {
		return errors.New("consumer is closed")
	}
	rb.readers.Add(1)
	defer rb.readers.Done()

	err := rb.consumer.Subscribe(topic, rebalanceCallback)
	if err != nil {
		return err
	}

	for !rb.stop {
		msg, err := rb.consumer.ReadMessage(ReadTimeout) // read message with timeout
		if err == nil {
			fmt.Printf("Message on topic %s, partition %d: key = %s, value = %s\n", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, string(msg.Key), string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() { // TODO process timeout
			slog.Error("Consumer error: ", slog.String("error", err.Error()))
		}
	}
	return nil
}

// Close stops reading messages
func (rb *Rebalanceconsumer) Close() {
	if !rb.stop {
		rb.stop = true // command to stop reading messages
		rb.readers.Wait()

		err := rb.consumer.Close()
		if err != nil {
			slog.Error("Failed to close consumer: ", slog.String("error", err.Error()))
		}
	}
}

// New creates a new consumer
func New(brokers string, group string) *Rebalanceconsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":             brokers,
		"group.id":                      group,
		"session.timeout.ms":            defaultSessionTimeout,
		"enable.auto.commit":            "true",
		"auto.commit.interval.ms":       1000,
		"auto.offset.reset":             "earliest", // earliest - сообщения с commit offset; latest - новые сообщение
		"partition.assignment.strategy": "cooperative-sticky",
	})
	if err != nil {
		return nil
	}

	return &Rebalanceconsumer{
		consumer: consumer,
		stop:     false,
	}
}

// rebalanceCallback
// Вызывается при ребалансировке
func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {

	switch ev := event.(type) {
	case kafka.AssignedPartitions: // Вызывается при назначении разделов
		fmt.Printf("AssignedPartitions, consumer: %s, protocol: %s. ", c.String(), c.GetRebalanceProtocol())
		fmt.Printf("Partitions: %v\n", ev.Partitions)

		// какие-то подготовительные действия

	case kafka.RevokedPartitions:
		fmt.Printf("RevokedPartitions, consumer: %s, protocol: %s. ", c.String(), c.GetRebalanceProtocol())
		fmt.Printf("Partitions: %v\n", ev.Partitions)

		// Зафиксировать смещения последних обработанных событий
		// освободить ресурсы если необходимо
	}

	return nil
}
