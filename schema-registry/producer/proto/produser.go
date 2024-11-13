package proto

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"kafka_examples/schema-registry/model"
	pb "kafka_examples/schema-registry/producer/proto/pkg/kafkaExample.v1"
	"log/slog"
	"sync"
)

const nullOffset = -1

type protoProducer struct {
	producer   *kafka.Producer
	serializer serde.Serializer
	stop       bool           // command to close
	senders    sync.WaitGroup // runned senders
}

func toProto(msg model.Human) *pb.ExampleMessage {
	return &pb.ExampleMessage{
		Id:      msg.Id,
		Name:    msg.Name,
		Age:     msg.Age,
		Student: msg.Student,
		City:    msg.City,
	}
}

// ProduceMessage
// Асинхронная отправка сообщения с ожиданием ответа.
// При возникновении ошибки offset = -1
func (sr *protoProducer) ProduceMessage(key string, msg model.Human, topic string) (offset int64, err error) {
	if sr.stop {
		return nullOffset, errors.New("producer is closed")
	}
	sr.senders.Add(1)
	defer sr.senders.Done()

	payload, err := sr.serializer.Serialize(topic, toProto(msg))
	if err != nil {
		return nullOffset, err
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = sr.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            []byte(key),
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		return nullOffset, err
	}

	event := <-deliveryChan
	switch ev := event.(type) {
	case *kafka.Message:
		return int64(ev.TopicPartition.Offset), nil
	case kafka.Error:
		return nullOffset, ev
	}

	return nullOffset, errors.New(event.String()) // не должно быть, в отчете о доставке пришло что-то неожиданное
}

func (sr *protoProducer) Close() {
	sr.stop = true    // command to close
	sr.senders.Wait() //waiting for the completion of the launched message sending

	err := sr.serializer.Close()
	if err != nil {
		slog.Error("Failed to close serializer: ", slog.String("error", err.Error()))
	}
	sr.producer.Close()
}

func NewProducer(brokers string, srURL string) (*protoProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(srURL))
	if err != nil {
		return nil, err
	}

	protocfg := protobuf.SerializerConfig{
		SerializerConfig: serde.SerializerConfig{},
		CacheSchemas:     true,
	}
	protocfg.AutoRegisterSchemas = true

	serializer, err := protobuf.NewSerializer(
		sr,
		serde.ValueSerde,
		&protocfg)
	if err != nil {
		return nil, err
	}

	return &protoProducer{
		producer:   producer,
		serializer: serializer,
		stop:       false,
	}, nil
}
