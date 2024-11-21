# Батчинг + сжатие
## Пакетная передача сообщений

Ключом для достижения высокой пропускной способности является отправка сообщений батчами(пакетами) - накопление сообщений в локальной очереди перед отправкой их одним батчем в kafka.

Накладные расход на запись постоянны и не поддаются уменьшению. На каждое сообщение 61 байт.

### Producer
Для включения пакетной передачи нужно при создании продюсера указать параметр
```
"go.batch.producer": true

"linger.ms": 5 (default) время накопления сообщений перед созданием батча и передачей его брокерам (0..900_000)
"batch.num.messages": 10_000 (default) максимальное количество сообщений в батче
"batch.size": 1_000_000 (default) максимальный размер батча в байтах

"message.max.bytes": 1_000_000 (default) максимальный размер запроса в Kafka
```

Чем больше пакет, тем выше пропускная способность.
Значение по умолчанию `"linger.ms": 5` не подходит для высокой пропускной способности, рекомендуется установить это значение >50 мс, при этом пропускная способность стабилизируется где-то в пределах 100–1000 мс в зависимости от шаблона и размера сообщений.

Для низкой задержки можно установить `linger.ms=0`. Эффективность пакетной передачи соответственно будет минимальна.

В confluent-kafka-go пакетная отправка в статусе EXPERIMENTAL.
Батчи никак не связаны с обычными сообщениями, они не поддерживают Headers и timestamps.
Соответственно невозможно отправить никакую дополнительную информацию по батчу, например контекст трассировки.

### Consumer
После создания консюмера librdkafka пытается сохранить `queued.min.messages` сообщений в локальной очереди.
librdkafka извлечет все потребляемые разделы, для которых этот брокер является лидером, с помощью одного запроса.

Ограничение настраивается в параметре:
```
"fetch.max.bytes": 52_428_800 (default)
`queued.min.messages` (default 100_000) минимальное количество сообщений, хранящееся в локальной очереди.
```

`fetch.max.bytes` должен быть >= `message.max.bytes` (default = 1_000_000).

Максимальный размер пакета сообщений принимаемый брокером определяется через `message.max.bytes` конфигурации брокера или `max.message.bytes` (конфигурации топика).

## Сжатие
Сжатие выполняется для пакета сообщений в локальной очереди.
Намного эффективнее сжимать батчи, чем отдельные сообщения.

Брокер распаковывает пакет для его проверки.
Например, брокер проверяет, что пакет содержит то же количество записей, что указано в заголовке пакета.
После проверки `пакет сообщений записывается на диск в сжатом виде`.
Пакет остается сжатым в журнале и передается потребителю в сжатом виде.
Потребитель распаковывает любые сжатые данные, которые он получает.

Kafka позволяет сжимать сообщения размером до 1 Гб, рекомендуемое значение - меньше 1Гб.

Сжимать зашифрованные сообщения нет смысла, коэффициент сжатия будет очень низким.

Сжатие может быть установлено как на продюсере, так и в свойствах топика.
Тип сжатия установленный в свойствах топика `compression.type` переопределяет тип сжатия продюсера с точки зрения хранения данных в кластере и отправке их потребителям.

### со стороны продюсера
`best practic` - нагрузка на сеть понижается, cpu увеличивается.

В свойствах топика должно быть установлено `compression.type: producer` (default).

Для включения сжатия нужно установить свойства в продюсере:
```
"compression.codec": "lz4"
"compression.level": -1 (default)
```

`compression.codec`: 
- none
- gzip `не рекомендуется`, высокие накладные расходы на кластер
- snappy
- lz4 `рекомендуется для высокой производительности` 
- zstd

`comporession.level`: Более высокое значение приведет к лучшему сжатию за счет большей нагрузки CPU.
- 0..9 --- для gzip
- 0..12 -- для lz4
- 0 ------ для snappy
- -1 ----- уровень сжатия для кодека по умолчанию.



### со стороны брокера
нагрузка на сеть от продюсера не меняется, cpu увеличивается.

В свойствах топика должно быть установлено `compression.type:` gzip, snappy, lz4, zstd, uncompressed.

Топики сжимаются на брокере. Если тип сжатия на продюсере не соответствует типу сжатия в свойствах топика, сообщения будут пережиматься в соответствии с типом указанным в свойствах топика.

Не рекомендуется.

### со стороны консюмера
Консюмер прозрачно обрабатывает как сжатые, так и несжатые сообщения.
Сжатие увеличивает производительность консюмера.
Увеличивает потребление CPU.

## Best practices

Батчинг + сжатие на продюсере.

## Documentation
[Apache Kafka Message Compression](https://www.confluent.io/blog/apache-kafka-message-compression/)

[Confluent Batch Processing for Efficiency](https://docs.confluent.io/kafka/design/efficient-design.html)