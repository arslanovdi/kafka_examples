package avro

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"kafka_examples/schema-registry/model"
	"log/slog"
	"sync"
	"time"
)

const (
	defaultSessionTimeout = 6000
	ReadTimeout           = 1 * time.Second // Не ставить слишком большим, влияет на время ожидания при остановке программы
)

type srConsumer struct {
	consumer     *kafka.Consumer
	deserializer *avro.GenericDeserializer
	stop         bool           // command to stop reading messages
	readers      sync.WaitGroup // runned readers
}

func NewConsumer(brokers string, srURL string, group string) (*srConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       brokers,
		"group.id":                group,
		"session.timeout.ms":      defaultSessionTimeout,
		"enable.auto.commit":      "true",
		"auto.commit.interval.ms": 1000,
		"auto.offset.reset":       "earliest", // earliest - сообщения с commit offset; latest - новые сообщение
	})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(srURL))
	if err != nil {
		return nil, err
	}

	deserializercfg := avro.NewDeserializerConfig()
	deserializercfg.UseLatestVersion = true

	deserializer, err := avro.NewGenericDeserializer(sr, serde.ValueSerde, deserializercfg)
	if err != nil {
		return nil, err
	}

	return &srConsumer{
		consumer:     consumer,
		deserializer: deserializer,
		stop:         false,
	}, nil

}

func (sr *srConsumer) Run(topic string, handler func(key string, msg model.Human, offset int64)) error {
	if sr.stop {
		return errors.New("consumer is closed")
	}
	sr.readers.Add(1)
	defer sr.readers.Done()

	err := sr.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}

	for !sr.stop {
		kafkaMsg, err := sr.consumer.ReadMessage(ReadTimeout) // read message with timeout
		if err == nil {
			h := model.Human{}
			err = sr.deserializer.DeserializeInto(
				topic,
				kafkaMsg.Value,
				&h)
			if err != nil {
				return err
			} else {
				handler(string(kafkaMsg.Key), h, int64(kafkaMsg.TopicPartition.Offset))
			}

		} else if !err.(kafka.Error).IsTimeout() { // TODO process timeout
			return err
		}
	}
	return nil
}

func (sr *srConsumer) Close() {
	sr.stop = true // command to stop reading messages
	sr.readers.Wait()

	err := sr.deserializer.Close()
	if err != nil {
		slog.Error("Failed to close deserializer: ", slog.String("error", err.Error()))
	}
	err = sr.consumer.Close()
	if err != nil {
		slog.Error("Failed to close consumer: ", slog.String("error", err.Error()))
	}
}
