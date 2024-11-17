package consumer

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka_examples/exactly-once/kafka-kafka/model"
	"kafka_examples/exactly-once/kafka-kafka/producer"
	"log/slog"
	"sync"
	"time"
)

const (
	defaultSessionTimeout = 6000
	ReadTimeout           = 1 * time.Second // Не ставить слишком большим, влияет на время ожидания при остановке программы
)

type srConsumer struct {
	consumer *kafka.Consumer
	stop     bool           // command to stop reading messages
	readers  sync.WaitGroup // runned readers
}

type Options struct {
	key   string
	value kafka.ConfigValue
}

func WithAutoCommit() Options {
	return Options{
		key:   "enable.auto.commit",
		value: true,
	}
}

func New(brokers string, group string, opts ...Options) (*srConsumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           group,
		"session.timeout.ms": defaultSessionTimeout,
		"enable.auto.commit": "false",
		"auto.offset.reset":  "earliest", // earliest - сообщения с commit offset; latest - новые сообщение
		"isolation.level":    "read_committed",
	}

	for _, opt := range opts {
		err := cfg.SetKey(opt.key, opt.value)
		if err != nil {
			return nil, err
		}
	}

	consumer, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}

	return &srConsumer{
		consumer: consumer,
		stop:     false,
	}, nil

}

func (sr *srConsumer) Run(topic string, p model.Producer, handler func(p model.Producer, msg *kafka.Message, cmsg producer.CommitMsgData)) error {
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
		msg, err := sr.consumer.ReadMessage(ReadTimeout) // read message with timeout
		if err == nil {

			meta, err := sr.consumer.GetConsumerGroupMetadata()
			if err != nil {
				return err
			}

			tp, err := sr.consumer.Position([]kafka.TopicPartition{{Topic: &topic, Partition: msg.TopicPartition.Partition}})
			if err != nil {
				return err
			}

			handler(p, msg, producer.CommitMsgData{
				Tp:   tp,
				Meta: meta,
			})

			continue
		}
		if !err.(kafka.Error).IsTimeout() { // TODO process timeout
			return err
		}
	}
	return nil
}

func (sr *srConsumer) Close() {
	sr.stop = true // command to stop reading messages
	sr.readers.Wait()

	err := sr.consumer.Close()
	if err != nil {
		slog.Error("Failed to close consumer: ", slog.String("error", err.Error()))
	}
}
