## Сбор статистики(метрик) с кластера kafka.

Консюмер можно настроить на получение статистики.
Статистика будет прилетать в формате `JSON` сообщением типа `*kafka.Stats`.
Для этого при создании консюмера нужно указать параметр:
```go
"statistics.interval.ms": 5000
```

Обрабатывать статистику там же где и ошибки.
```go
ev := consumer.Poll(100)
switch ev := event.(type) {
	case *kafka.Message:
	    if ev.TopicPartition.Error != nil {
			// обработка ошибки
	    }
		// события
	case kafka.Error:
		// ошибки 
	case *kafka.Stats:
		// статистика
        var stats map[string]interface{}
        json.Unmarshal([]byte(e.String()), &stats)
        fmt.Printf("Stats: %v messages (%v bytes) messages consumed\n", stats["rxmsgs"], stats["rxmsg_bytes"])
	}
```

Можно создать пустой топик, чисто для того, чтобы получить статистику.

### Documentation
[Список полей со статистикой(метриками)](https://github.com/confluentinc/librdkafka/blob/master/STATISTICS.md)