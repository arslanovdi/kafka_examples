package generator

import (
	"fmt"
	"kafka_examples/exactly-once/kafka-kafka/model"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Generate добавляет в топик случайное число с шагом ticker
func Generate(producer model.Producer, topic string, key string, ticker time.Duration) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	tik := time.NewTicker(ticker)

loop:
	for {
		select {
		case <-stop:
			break loop
		case <-tik.C:
			offset, err := producer.Produce(topic, key, []byte(fmt.Sprintf("%d", rand.Int31())))
			if err != nil {
				slog.Error("Failed to produce message: ", slog.String("error", err.Error()))
			}

			slog.Info("Message produced", slog.Int64("offset", offset))
		}
	}
}
