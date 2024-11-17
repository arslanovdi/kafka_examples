package producer

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"sync"
)

var ErrProducerClosed = errors.New("producer is closed")

const nullOffset = -1

type Options struct {
	key   string
	value kafka.ConfigValue
}

func WithTransaction(tid string) Options {
	return Options{
		key:   "transactional.id",
		value: tid,
	}
}

type CommitMsgData struct {
	Tp   []kafka.TopicPartition
	Meta *kafka.ConsumerGroupMetadata
}

type producer struct {
	producer *kafka.Producer
	stop     bool           // command to close
	senders  sync.WaitGroup // runned senders
	close    chan bool
}

// TransactionProduce
// транзакционная отправка сообщения и фиксация смещения консюмера
func (sr *producer) TransactionProduce(topic string, key string, payload []byte, cmsg CommitMsgData) error {
	if sr.stop {
		return ErrProducerClosed
	}
	sr.senders.Add(1)
	defer sr.senders.Done()

	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          payload,
	}

	err := sr.producer.BeginTransaction()
	if err != nil {
		return err
	}

	err = sr.producer.Produce(&msg, nil) // обработка событий в транзакционном продюсере идет по каналу producer.Events().
	if err != nil {
		return err
	}

	err = sr.producer.SendOffsetsToTransaction(nil, cmsg.Tp, cmsg.Meta)
	if err != nil {
		return err
	}

	/*ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()*/
	err = sr.producer.CommitTransaction(nil)

	if err != nil {
		return err
	}

	return nil
}

// Produce
// Асинхронная отправка сообщения с ожиданием ответа.
func (sr *producer) Produce(topic string, key string, payload []byte) (offset int64, err error) {
	if sr.stop {
		return nullOffset, errors.New("producer is closed")
	}
	sr.senders.Add(1)
	defer sr.senders.Done()

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
		if ev.IsFatal() {
			// all further produce requests will fail in idempotence producer mode
			slog.Error("Fatal error, terminating", slog.String("error", ev.Error()))
			panic(ev.Error())
		}
		return nullOffset, ev
	default:
		return nullOffset, errors.New(event.String()) // не должно быть, в отчете о доставке пришло что-то неожиданное
	}
}

func (sr *producer) Close() {
	sr.stop = true // command to close
	select {
	case sr.close <- true: // stop events routine
	default:
	}
	sr.senders.Wait() // waiting for the completion of the launched message sending

	close(sr.close)
	sr.producer.Close()
}

func New(brokers string, opts ...Options) (*producer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":                     brokers,
		"enable.idempotence":                    true,
		`max.in.flight.requests.per.connection`: 3,
		`retries`:                               3,
		"acks":                                  "all",
		//"transactional.id":                      "my-transactional-id",
	}

	withtransaction := false

	for _, opt := range opts {
		err := cfg.SetKey(opt.key, opt.value)
		if opt.key == "transactional.id" {
			withtransaction = true
		}
		if err != nil {
			return nil, err
		}
	}

	p, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}

	cl := make(chan bool)

	ip := producer{
		producer: p,
		stop:     false,
		close:    cl,
	}

	if withtransaction {
		// Обработка ошибок в отдельной горутине обязательна, т.к. транзакционные методы лочат вызывающие потоки.
		go func() { // events routine
			ip.senders.Add(1)
			defer ip.senders.Done()

			for {
				select {
				case event := <-p.Events():
					switch ev := event.(type) {
					case *kafka.Message:
						slog.Info("Message produced to topic", slog.String("topic", *ev.TopicPartition.Topic), slog.Int64("offset", int64(ev.TopicPartition.Offset)))
					case kafka.Error:
						if ev.IsFatal() {
							// all further produce requests will fail in idempotence producer mode
							slog.Error("Fatal error, terminating", slog.String("error", ev.Error()))
							panic(ev.Error())
						}
					}
				case <-ip.close:
					return
				}
			}
		}()

		err = p.InitTransactions(nil) // if use context without deadline, deadline =  2*transaction.timeout.ms
		if err != nil {
			ip.close <- true
			return nil, err
		}
	}

	return &ip, nil
}

func (sr *producer) Abort() {
	err := sr.producer.AbortTransaction(context.Background())
	if err != nil {
		slog.Error("Failed to abort transaction: ", slog.String("error", err.Error()))
	}
}
