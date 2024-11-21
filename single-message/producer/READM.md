# Producer

- Продюсер создается при помощи `kafka.NewProducer()`, с указанием как минимум свойства `bootstrap.servers`.
```go
producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
```

- Сообщения типа `*kafka.Message` отправляются на канал `.ProduceChannel`, либо вызовом `.Produce()`.
```go
err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            []byte(key),
		Value:          []byte(value),
	}, nil)
```

Отправка сообщения является асинхронной операцией. Уведомление о доставке сообщения производится при помощи отчетов о доставке.

Отчет прилетают на канал `producer.Events()`, или на канал, переданный параметром при вызове `.Produce()`.
```go
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            []byte(key),
		Value:          []byte(value),
	}, deliveryChan)
```

На тот же канал прилетают ошибки типа `kafka.Error`.

Если отчеты о доставке не нужны, их можно полностью отключить, установив свойство конфигурации `"go.delivery.reports": false`.