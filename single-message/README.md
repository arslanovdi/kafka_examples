# Семантики
Kafka гарантирует доставку at-least-once по умолчанию и позволяет пользователю реализовать доставку at-most-once,
отключив повторные попытки на производителе и зафиксировав смещения в потребителе перед обработкой пакета сообщений.

## at most once (максимум один раз)
### producer
Продюсер считает сообщение отправленным без ожидания ответа от брокера, нет ретраев.

### consumer
1. читать сообщение
2. сохранять позицию(offset)
3. обрабатывать сообщение

Существует вероятность, что процесс потребитель даст сбой после сохранения своей позиции, но до завершения обработки сообщения.

Новый потребитель начнет читать сообщения с сохраненной позиции, соответственно сообщение которое не успел обработать предыдущий потребитель будет потеряно.

## at least once (минимум один раз)
### producer
Продюсер ждет подтверждения от брокера. Если ответа нет, то оно будет отправлено еще раз.
Возможны дубли (если продюсер не идемпотентный).

### consumer
1. читать сообщение
2. обрабатывать сообщения
3. сохранять позицию(offset)

Существует вероятность, что процесс потребитель даст сбой после обработки сообщения, но до сохранения своей позиции.

Новый потребитель получит сообщения которые уже были обработаны предыдущим потребителем.

## exactly once (точно один раз)
Описание в соответствующем примере.

