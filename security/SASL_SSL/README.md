## Авторизация
Для кластеров kraft авторизация поставляется классом `org.apache.kafka.metadata.authorizer.StandardAuthorizer` для обработки ACL (создание, чтение, запись, удаление).

## SASL методы аутентификации
- [GSSAPI (Kerberos)](https://kafka.apache.org/documentation/#security_sasl_kerberos)
- [PLAIN](https://kafka.apache.org/documentation/#security_sasl_plain)
- [SCRAM-SHA-256](https://kafka.apache.org/documentation/#security_sasl_scram)
- [SCRAM-SHA-512](https://kafka.apache.org/documentation/#security_sasl_scram)
- [OAUTHBEARER](https://kafka.apache.org/documentation/#security_sasl_oauthbearer)

## Example
В примере используется аутентификация с использованием SASL/PLAIN, авторизация с использованием ACL и SSL шифрованием.
SASL/PLAIN включен для EXTERNAL Listener.

### Конфигурация брокера
```
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:SSL,INTERNAL:SSL,EXTERNAL:SASL_SSL   # список листенеров с протоколами доступа
KAFKA_LISTENER_NAME_INTERNAL_SSL_CLIENT_AUTH: required
KAFKA_LISTENER_NAME_EXTERNAL_SSL_CLIENT_AUTH: required
KAFKA_LISTENER_NAME_CONTROLLER_SSL_CLIENT_AUTH: required
KAFKA_ALLOW_EVERYONE_IF_NO_ACL_FOUND: false
KAFKA_SUPER_USERS: User:kafka1;User:kafka2;User:kafka3;User:kafkaui;
KAFKA_SSL_PRINCIPAL_MAPPING_RULES: RULE:^.*[Cc][Nn]=([a-zA-Z0-9.]*).*$/$1/L,,DEFAULT
KAFKA_AUTHORIZER_CLASS_NAME: org.apache.kafka.metadata.authorizer.StandardAuthorizer
KAFKA_LISTENER_NAME_EXTERNAL_PLAIN_SASL_JAAS_CONFIG: org.apache.kafka.common.security.plain.PlainLoginModule required user_producer="producer-secret" user_consumer="consumer-secret";
KAFKA_SASL_ENABLED_MECHANISMS: PLAIN
```

#### Немного подробнее про параметры, тут есть подводные камни:
Параметр `KAFKA_AUTHORIZER_CLASS_NAME:`
Для кластера Kraft параметр должен быть равен `org.apache.kafka.metadata.authorizer.StandardAuthorizer`

После добавления этого параметра кластер падает с ошибкой ACL. Это происходит из-за того, что ACL начинает проверять все листенеры, включая контроллеры. А уровни доступа еще не заданы.

В списках контроля доступа(ACL) имя пользователя(Principal) зависит от протокола.

SSL Principal - это строка +- в формате `"CN=writeuser,OU=Unknown,O=Unknown,L=Unknown,ST=Unknown,C=Unknown"`.

Для удобства вырезаем из строки имя хоста, используя маппинг, который задается параметром:

`KAFKA_SSL_PRINCIPAL_MAPPING_RULES: RULE:^.*[Cc][Nn]=([a-zA-Z0-9.]*).*$/$1/L,,DEFAULT`

Получаем kafka1;kafka2;kafka3;kafkaui;

Этих пользователей указываем как SuperUser, после этого кластер стартует без ошибок, т.к. права выданы.
`KAFKA_SUPER_USERS: User:kafka1;User:kafka2;User:kafka3;User:kafkaui;`

SASL/PLAIN Principal - username, которые задаются в одном из параметров
- `KAFKA_LISTENER_NAME_EXTERNAL_PLAIN_SASL_JAAS_CONFIG`: org.apache.kafka.common.security.plain.PlainLoginModule required user_producer="producer-secret" user_consumer="consumer-secret";
- `KAFKA_OPTS`: "-Djava.security.auth.login.config=/etc/kafka/jaas/kafka3_jaas.conf"

В kafka3 реализован вариант с `KAFKA_OPTS`.

### Конфигурация клиента (продюсер/консюмер)

```
"security.protocol": "sasl_ssl",
"sasl.mechanism":    "PLAIN",
"sasl.username":     "username",
"sasl.password":     "secret",
```

### Настройка ACL
Настроить ACL можно через API Kafka AdminClient или при помощи CLI команды `kafka-acls.sh`.

[Kafka Authorization management CLI](https://kafka.apache.org/documentation/#security_authz_cli) можно найти в каталоге /opt/kafka/bin вместе со всеми остальными CLI.

kakfa-acls.sh является клиентом по отношению к брокеру. Для установки шифрованного соединение, нужно передать параметры подключения через --command-config.

Подключаюсь к INTERNAL listener при помощи сертификата, который лежит на kafka1. )))
На kafka1 проброшен файл конфигурации клиента `client.properties`.

Пользователю `producer` даем права `create/describe/write` на топик `topic` командой:

```bash
docker exec kafka1 /opt/kafka/bin/kafka-acls.sh --command-config /client.properties --bootstrap-server kafka1:9091,kafka2:9091,kafka3:9091 --add --allow-principal User:producer --producer --topic topic
```

`При подпытки записи без прав доступа librdkafka не выдает ошибку, просто возвращает offset=-1001 !`

Пользователю `consumer` даем права `read` на все топики и группы, `describe` на все топики, командой:

```bash
docker exec kafka1 /opt/kafka/bin/kafka-acls.sh --command-config /client.properties --bootstrap-server kafka1:9091,kafka2:9091,kafka3:9091 --add --allow-principal User:consumer --consumer --topic "*" --group "*"
```

Kafka-UI умеет удалять ACL, а создавать похоже нет. :)

## Documentation
[Apache Kafka Authentication using SASL](https://kafka.apache.org/documentation/#security_sasl)

[Apache Kafka Authorization and ACLs](https://kafka.apache.org/documentation/#security_authz)

[Use SASL/SCRAM authentication in Confluent Platform](https://docs.confluent.io/platform/current/security/authentication/sasl/scram/overview.html#auth-sasl-scram-broker-config)

[librdkafka configuration](https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md)