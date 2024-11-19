# Ребалансировка партиций.
Изменение количества консюмеров приводит к ребалансировке партиций.

## Стратегии ребалансировки
Для изменения стратегии ребалансировки необходимо установить свойство конфигурации консюмера
```
"partition.assignment.strategy"
```

Список стратегий Librdkafka:
- `"range"`
- `"roundrobin"`
- `"cooperative-sticky"`

Значение по умолчанию в librdkafka:
`"range,roundrobun"`

Лидер консюмер группы будет использовать стратегию, поддерживаемую всеми участниками.
Если поддерживается несколько стратегий использоваться будет та, которая идет первой в списке.

### range
`Stop The World`
На время перебалансировки чтение топиков останавливается.

Получение партиций со всех топиков, на которые подписаны.

Например 3 топика (TA, TB, TC) по 3 партиции (P0,P1,P2).
4 кнсюмера получат следующие партиции:
1. TA-P0, TB-P0, TC-P0
2. TA-P1, TB-P1, TC-P1
3. TA-P2, TB-P2, TC-P2
4. ничего не получит ))
### roundrobin
`Stop The World`
На время перебалансировки чтение топиков останавливается.

Равномерное распределение партиций между консюмерами.

Например 3 топика (TA, TB, TC) по 3 партиции (P0,P1,P2).
4 кнсюмера получат следующие партиции:
1. TA-P0, TB-P1, TC-P2
2. TA-P1, TB-P2
3. TA-P2, TC-P0
4. TB-P0, TC-P1


### cooperative-sticky
Без `Stop The World`.
[KIP-429](https://cwiki.apache.org/confluence/display/KAFKA/KIP-429%3A+Kafka+Consumer+Incremental+Rebalance+Protocol)

Как roundrobin, но
затрагиваются только те партиции, которым нужно сменить назначение.

## Обработка событий
Для обработки событий ребалансировки необходимо добавить обработчик событий при подписке на топик.

```go
_ = consumer.Subscribe("test", rebalanceCallback)

func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {
    switch ev := event.(type) {
    case kafka.AssignedPartitions:
        fmt.Println(c.GetRebalanceProtocol())
        fmt.Println(ev.Partitions)

    case kafka.RevokedPartitions:
        fmt.Println(c.GetRebalanceProtocol())
        fmt.Println(ev.Partitions)
    }
return nil
}
```

Сообщения прилетают при каждой ребалансировке, если консюмер не затронут ev.Partitions будет пуст [].

### `kafka.AssignedPartitions`
Вызывается при назначении партиций, до того как консюмер начнет получать сообщения.
Тут можно выполнить подготовительные операции.

Любая подготовка должна гарантированно выполниться в течении `max.poll.interval.ms` (default: 300000).

### `kafka.RevokedPartitions`
Вызывается когда консюмер должен отказаться от партиций, которыми ранее владел. При перебалансировке или закрытии консюмера.

Тут нужно зафиксировать смещения, чтобы новый консюмер начал получать сообщения с нужных смещений.

## [Next generation rebalance protocol](https://github.com/confluentinc/librdkafka/blob/master/INTRODUCTION.md#next-generation-of-the-consumer-group-protocol-kip-848)
Статус - Experimental.

Начиная с версии librdkafka [2.4.0](https://github.com/confluentinc/librdkafka/releases/tag/v2.4.0) представлен протокол ребалансировки групп потребителей следующего поколения, определенный в KIP 848. Пока в раннем доступе.

При использовании этого протокола роль лидера группы (участника) устраняется, а назначение рассчитывается координатором группы (брокером) и отправляется каждому участнику посредством тактовых импульсов.

Для его тестирования необходимо настроить кластер Kafka в режиме KRaft и включить новый групповой протокол с помощью свойства `group.coordinator.rebalance.protocols`.
Версия брокера должна быть Apache Kafka 3.7.0 или более поздняя. См. Apache Kafka [Release Notes](https://cwiki.apache.org/confluence/display/KAFKA/The+Next+Generation+of+the+Consumer+Rebalance+Protocol+%28KIP-848%29+-+Early+Access+Release+Notes).