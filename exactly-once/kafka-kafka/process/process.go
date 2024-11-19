package process

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kafka_examples/exactly-once/kafka-kafka/model"
	"kafka_examples/exactly-once/kafka-kafka/producer"
	"log/slog"
	"strconv"
)

var newTopic = "new-topic"

// Multiplier
// умножает значение msg на 10, отправляет его в newTopic и фиксирует смещение
func Multiplier(p model.Producer, msg *kafka.Message, cmsg producer.CommitMsgData) {
	i, _ := strconv.Atoi(string(msg.Value))
	i *= 10
	value := []byte(fmt.Sprintf("%d", i))

	// Пытаемся выполнить транзакцию, пока не получим успех или Фатальную ошибку
loop:
	for {
		err := p.TransactionProduce(newTopic, string(msg.Key), value, cmsg)
		switch ev := err.(type) {
		case nil: // commit transaction прошел успешно
			return
		case kafka.Error:
			if ev.TxnRequiresAbort() {
				slog.Error("Abort transaction", slog.String("error", ev.Error()))
				p.Abort() //do_abort_transaction_and_reset_inputs()
				continue
			}
			if ev.IsRetriable() {
				slog.Error("Retry transaction", slog.String("error", ev.Error()))
				continue
			}
			// treat all other errors as Fatal errors
			slog.Error("Fatal transaction error: ", slog.String("error", err.Error()))
			panic("process.Multiplier " + err.Error())
		default:
			if errors.Is(err, producer.ErrProducerClosed) { // выходим если продюсер закрыт
				slog.Error("producer is closed")
				break loop
			}
			slog.Error("transaction error: ", slog.String("error", err.Error())) // неизвестная ошибка
			break loop
		}
	}
}
