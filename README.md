## Apache Kafka - платформа потоковой передачи событий.

Для использования библиотеки [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) должен быть установлен
```go
CGO_ENABLED=1
```

## Examples

### single messages
одиночные сообщения producer/consumer

### batch messages
батчи producer/consumer // EXPERIMENTAL

### schema registry
[Пример producer/consumer с использованием schema registry, сериализацией в форматах protobuf/avro/json](https://github.com/arslanovdi/kafka_examples/tree/master/schema-registry)

### [exactly once](https://github.com/arslanovdi/kafka_examples/tree/master/exactly-once)
- [kafka-kafka](https://github.com/arslanovdi/kafka_examples/tree/master/exactly-once/kafka-kafka)
- kafka-postgres

### kafka sreams
использование kafka sreams

### authentication
пример аутентификации

### partitioner
пример своей реализации разбиения ключей по партициям

