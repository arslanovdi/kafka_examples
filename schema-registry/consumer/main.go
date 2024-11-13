package main

import (
	"fmt"
	"kafka_examples/schema-registry/consumer/avro"
	"kafka_examples/schema-registry/consumer/json"
	"kafka_examples/schema-registry/consumer/proto"
	"kafka_examples/schema-registry/model"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	brokers        = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"
	schemaRegistry = "http://127.0.0.1:8081"

	prototopic = "protobuf-topic"
	avrotopic  = "avro-topic"
	jsontopic  = "json-topic"
)

type SRConsumer interface {
	Run(topic string, handler func(key string, msg model.Human, offset int64)) error
	Close()
}

func handleProtoMessage(key string, message model.Human, offset int64) {
	fmt.Printf("protobuf message\nkey:%s\nmessage:%v\noffset:%d\n\n", key, message, offset)
}
func handleAvroMessage(key string, message model.Human, offset int64) {
	fmt.Printf("avro message\nkey:%s\nmessage:%v\noffset:%d\n\n", key, message, offset)
}
func handleJSONMessage(key string, message model.Human, offset int64) {
	fmt.Printf("json message\nkey:%s\nmessage:%v\noffset:%d\n\n", key, message, offset)
}

func main() {

	// Protobuf
	pconsumer, err := proto.NewConsumer(brokers, schemaRegistry, "myProtoGroup")
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pconsumer.Close()
	go func() {
		err = pconsumer.Run(prototopic, handleProtoMessage)
		if err != nil {
			slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Avro
	aconsumer, err := avro.NewConsumer(brokers, schemaRegistry, "myAvroGroup")
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer aconsumer.Close()
	go func() {
		err = aconsumer.Run(avrotopic, handleAvroMessage)
		if err != nil {
			slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// JSON
	jconsumer, err := json.NewConsumer(brokers, schemaRegistry, "myJSONGroup")
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer jconsumer.Close()
	go func() {
		err = jconsumer.Run(jsontopic, handleJSONMessage)
		if err != nil {
			slog.Error("Failed to run consumer: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
