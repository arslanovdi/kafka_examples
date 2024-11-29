## Apache Kafka - платформа потоковой передачи событий.
Для использования библиотеки [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) должен быть установлен
```go
CGO_ENABLED=1
```

## Examples

### single messages
[Пример одиночные сообщения producer/consumer](https://github.com/arslanovdi/kafka_examples/tree/master/single-message)

### Батчинг и сжатие сообщений
[Пример](https://github.com/arslanovdi/kafka_examples/tree/master/batch_messages)

### schema registry
[Пример producer/consumer с использованием schema registry, сериализацией в форматах protobuf/avro/json](https://github.com/arslanovdi/kafka_examples/tree/master/schema-registry)

### [exactly once](https://github.com/arslanovdi/kafka_examples/tree/master/exactly-once)
- [kafka-kafka](https://github.com/arslanovdi/kafka_examples/tree/master/exactly-once/kafka-kafka)
- kafka-postgres

### kafka streams
использование kafka streams

### security
[Примеры](https://github.com/arslanovdi/kafka_examples/tree/master/security)

### rebalance
[Пример обработки событий ребалансировки партиций](https://github.com/arslanovdi/kafka_examples/tree/master/rebalance)

### librdkafka statistics from consumer events
[Получаем метрики kafka кластера через события консюмера](https://github.com/arslanovdi/kafka_examples/tree/master/stats)

## Documentation
[Apache kafka](https://kafka.apache.org/documentation/)

[Kafka Improvement Proposals](https://cwiki.apache.org/confluence/display/KAFKA/Kafka+Improvement+Proposals)

[Introduction to librdkafka](https://github.com/confluentinc/librdkafka/blob/master/INTRODUCTION.md)

[librdkafka configuration properties](https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md)

[Confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)

[docs.confluent.io](https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html)