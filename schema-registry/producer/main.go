package main

import (
	"kafka_examples/schema-registry/model"
	"kafka_examples/schema-registry/producer/avro"
	"kafka_examples/schema-registry/producer/json"
	"kafka_examples/schema-registry/producer/proto"
	"log/slog"
	"os"
)

const (
	brokers        = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	schemaRegistry = "http://127.0.0.1:8081"

	prototopic = "protobuf-topic"
	avrotopic  = "avro-topic"
	jsontopic  = "json-topic"
)

type SRProducer interface {
	ProduceMessage(key string, msg model.Human, topic string) (int64, error)
	Close()
}

func main() {

	// Protobuf
	pproducer, err := proto.NewProducer(brokers, schemaRegistry)
	if err != nil {
		slog.Error("Failed to create proto producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pproducer.Close()

	sendMessages(pproducer, prototopic)

	// Avro
	aproducer, err := avro.NewProducer(brokers, schemaRegistry)
	if err != nil {
		slog.Error("Failed to create avro producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer aproducer.Close()

	sendMessages(aproducer, avrotopic)

	// JSON
	jproducer, err := json.NewProducer(brokers, schemaRegistry)
	if err != nil {
		slog.Error("Failed to create json producer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer jproducer.Close()

	sendMessages(jproducer, jsontopic)
}

// sendMessages
// отправка тестовых сообщений
func sendMessages(producer SRProducer, topic string) {
	msg := model.Human{
		Id:      0,
		Name:    "Ivan",
		Age:     22,
		Student: true,
		City:    nil,
	}
	offset, err := producer.ProduceMessage("without city", msg, topic)
	if err != nil {
		slog.Error("Failed to produce message: ", slog.String("error", err.Error()))
	} else {
		slog.Info("Message produced to topic", slog.String("topic", topic), slog.Int64("offset", offset))
	}

	city := "Moscow"
	msg = model.Human{
		Id:      0,
		Name:    "Ivan",
		Age:     22,
		Student: true,
		City:    &city,
	}
	offset, err = producer.ProduceMessage("Moscow", msg, topic)
	if err != nil {
		slog.Error("Failed to produce message: ", slog.String("error", err.Error()))
	} else {
		slog.Info("Message produced to topic", slog.String("topic", topic), slog.Int64("offset", offset))
	}
}
