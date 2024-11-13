package avro

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
	"kafka_examples/schema-registry/model"
	"log/slog"
	"sync"
)

const nullOffset = -1

type avroProducer struct {
	producer   *kafka.Producer
	serializer serde.Serializer
	stop       bool           // command to close
	senders    sync.WaitGroup // runned senders
}

// ProduceMessage
// Асинхронная отправка сообщения с ожиданием ответа.
// При возникновении ошибки offset = -1
func (sr *avroProducer) ProduceMessage(key string, msg model.Human, topic string) (offset int64, err error) {
	if sr.stop {
		return nullOffset, errors.New("producer is closed")
	}
	sr.senders.Add(1)
	defer sr.senders.Done()

	payload, err := sr.serializer.Serialize(topic, &msg)
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

func (sr *avroProducer) Close() {
	sr.stop = true    // command to close
	sr.senders.Wait() // waiting for the completion of the launched message sending

	err := sr.serializer.Close()
	if err != nil {
		slog.Error("Failed to close serializer: ", slog.String("error", err.Error()))
	}
	sr.producer.Close()
}

func NewProducer(brokers string, srURL string) (*avroProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(srURL))
	if err != nil {
		return nil, err
	}

	avrocfg := avrov2.NewSerializerConfig()
	avrocfg.AutoRegisterSchemas = true

	serializer, err := avrov2.NewSerializer(
		sr,
		serde.ValueSerde,
		avrocfg)
	if err != nil {
		return nil, err
	}

	return &avroProducer{
		producer:   producer,
		serializer: serializer,
		stop:       false,
	}, nil
}
