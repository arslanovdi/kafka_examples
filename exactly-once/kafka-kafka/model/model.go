package model

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka_examples/exactly-once/kafka-kafka/producer"
)

type Consumer interface {
	Run(topic string, p Producer, handler func(p Producer, msg *kafka.Message, cmsg producer.CommitMsgData)) error
	Close()
}

type Producer interface {
	TransactionProduce(topic string, key string, payload []byte, cmsg producer.CommitMsgData) error
	Produce(topic string, key string, payload []byte) (offset int64, err error)
	Close()
	Abort()
}
