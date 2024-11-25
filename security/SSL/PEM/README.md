## Overview
Начиная с Apacha Kafka 2.7 можно использовать сертификаты TLS в формате PEM.

PEM — это схема кодирования сертификатов x509 и закрытых ключей в виде строк Base64 ASCII.
Это упрощает обработку ваших сертификатов. 
Вы можете просто предоставить ключи и сертификаты приложению в виде строковых параметров (например, через переменные среды). 
`Это особенно полезно, если ваши приложения работают в контейнерах, где монтирование файлов в контейнеры немного усложняет конвейер развертывания.`

Вы можете добавлять сертификаты непосредственно в файл конфигурации ваших клиентов или брокеров.
Если вы предоставляете их в виде однострочных строк, вы должны преобразовать исходный многострочный формат в одну строку, добавив символы перевода строки ( \n ) в конце каждой строки.

## Example step by step

Все команды можно посмотреть в makefile

Все ключи и сертификаты создаются при помощи openssl, без шифрования. Незашифрованные ключи допустимы при передаче ключей в конфигурацию в виде строк, и хранении их в защищенном хранилище, Vault например.

Kafka SSL engine factory по умолчанию поддерживает формат PEM с ключами PKCS#8.
Поэтому производится преобразование ключей в формат PKCS#8.
### 1. Создаем корневой ключ и корневой сертификат CA
```bash
make createCA
```

### 2. Создаем PEM сертификаты для каждого брокера и клиентов
```bash
make createpem
```
Команда выполняет поочередно:
- 2.1 - создает ключ и запрос на получение сертификата (.key и .csr)
- 2.2 - преобразует ключ в формат PKCS#8
- 2.3 - подписывает .csr при помощи root CA. С добавлением в сертификат SAN с dns именами.
- 2.4 - сохраняет приватный ключ и подписанный сертификат в один файл (.pem)

Про пункт 2.3 немного подробнее:

Указанные в SAN DNS имена используются Kafka для проверки соответствия сертификата хосту. В примере в SAN добавляется 2 DNS имени - имя брокера и localhost.

localhost добавлен исключительно для тестирования, т.к. с локальной машины не резолвятся хостнеймы внутри докер контейнера. Это же имя указано в конфигурации брокеров "EXTERNAL://localhost:2909x"
В проде вместо localhost Должно быть имя конкретной машины.

В проде сертификат выдается для конкретной машины, и в SAN DNS должно быть только имя этой машины.

Либо сертификат на пул машин..., тогда соответственно в SAN список всех машин.

### 3. Сконфигурировать брокеры и клиенты
```bash
make config
```
Команда формирует конфигурации для producer и consumer. Их нужно указать в kafka.ConfigMap при создании producer и consumer.

Для примера реализованы разные подходы по передачи ключа и сертификата.

- producer.properties - конфигурация с передачей ключа и сертификата в виде файлов
- consumer.properties - конфигурация с передачей ключа и сертификата в виде строки

```bash
make kafka3
```
- kafka3.properties - пример конфигурации брокера, которая будет принимать сертификаты в виде строки

скопировать в раздел environment - kafka3 - compose.yaml

Настройка брокера kafka1 с монтированием файлов сертификатов (в моем примере это kafka1, kafka2):
```
   environment:
      KAFKA_LISTENERS: INTERNAL://:9091,CONTROLLER://kafka1:9093,EXTERNAL://:29092           # kafka1 - должно быть в SAN сертификата контроллера
      KAFKA_ADVERTISED_LISTENERS: INTERNAL://kafka1:9091,EXTERNAL://localhost:29092          # localhost для тестирования, в проде должно быть имя конкретной машины
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:SSL,INTERNAL:SSL,EXTERNAL:SSL         # Включаем SSL для контроллера и INTERNAL и EXTERNAL клиентов
      # SSL use pem files
      KAFKA_SECURITY_PROTOCOL: SSL
      KAFKA_SSL_KEYSTORE_TYPE: PEM
      KAFKA_SSL_KEYSTORE_LOCATION: /etc/kafka/secrets/kafka1-keypair.pem
      KAFKA_SSL_TRUSTSTORE_TYPE: PEM
      KAFKA_SSL_TRUSTSTORE_LOCATION: /etc/kafka/secrets/rootCA.crt
      # SSL clients
      KAFKA_SSL_CLIENT_AUTH: required
   volumes:
      - type: bind
        source: ./certs/kafka1-keypair.pem
        target: /etc/kafka/secrets/kafka1-keypair.pem
      - type: bind
        source: ./certs/rootCA.crt
        target: /etc/kafka/secrets/rootCA.crt
```

### Kafka-UI
Судя по всему проверки соответствия имени хоста и поля SAN DNS сертификата не происходят.
Менял имя хоста kafkaui, все равно все работает.

Kafka-UI принимает CA сертификат только в формате jks.

Следующая команда загружает rootCA.crt в файл rootCA.jks, с паролем "qwerty".
```bash
make kafkaui
```

Конфигурация:
```
   # SSL
   KAFKA_CLUSTERS_0_PROPERTIES_SECURITY_PROTOCOL: SSL
   KAFKA_CLUSTERS_0_PROPERTIES_SSL_KEYSTORE_TYPE: PEM
   KAFKA_CLUSTERS_0_PROPERTIES_SSL_KEYSTORE_LOCATION: /kafkaui-keypair.pem
   KAFKA_CLUSTERS_0_PROPERTIES_SSL_TRUSTSTORE_LOCATION: /rootCA.jks
   KAFKA_CLUSTERS_0_PROPERTIES_SSL_TRUSTSTORE_PASSWORD: qwerty  
   volumes:
      - type: bind
        source: ./certs/kafkaui-keypair.pem
        target: /kafkaui-keypair.pem
      - type: bind
        source: ./certs/rootCA.jks
        target: /rootCA.jks
```

## Распространенные ошибки при настройке цепочки сертификатов
1. Если ваш закрытый ключ зашифрован, вам необходимо преобразовать его из формата PKCS#1 в формат PKCS#8, чтобы Java/Kafka могли его правильно прочитать.
2. Если вы хотите предоставить сертификат PEM в виде однострочной строки, обязательно добавьте символы перевода строки в конце каждой строки `\n`.
   В противном случае сертификат будет считаться недействительным.
3. Цепочка сертификатов должна включать ваш сертификат вместе со всеми промежуточными сертификатами CA, которые его подписали, в порядке подписания.
   Так, например, если ваш сертификат был подписан сертификатом A, который был подписан сертификатом B, который был подписан корневым сертификатом, ваша цепочка сертификатов должна включать:
   ваш сертификат, сертификат A и сертификат B, в том же порядке. Обратите внимание, что корневого сертификата не должно быть в цепочке.
4. Порядок сертификатов в вашей цепочке сертификатов важен (см. пункт 3)

### Documentation
[Kafka security documentation](https://kafka.apache.org/090/documentation.html#security)

[How to use PEM certificates with Apache Kafka](https://codingharbour.com/apache-kafka/using-pem-certificates-with-apache-kafka/)

[librdkafka configuration properties](https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md)