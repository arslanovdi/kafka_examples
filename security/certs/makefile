DNS2=Ubuntu

# Создаем корневой ключ и корневой сертификат CA
createroot:
	openssl req \
        -x509 \
        -days 365 \
        -newkey rsa:2048 \
        -keyout root.key \
        -out root.crt \
        -noenc \
        -subj "/CN=root.kafkaexample.ru/OU=TEST/O=KafkaExample/L=Moscow/C=RU"

# Создаем PEM сертификаты для каждого брокера и клиентов
SSLClients=kafka1 kafka2 kafka3 kafkaui producer consumer
createpem:
	for certname in $(SSLClients) ; do \
		make .createkeypair certname=$$certname ; \
		make .converttopkcs8 certname=$$certname ; \
		make .signcsr dns1=$$certname ; \
		make .combine certname=$$certname ; \
		rm *.srl ; \
	done

# создаем ключ и запрос на получение сертификата
.createkeypair:
	openssl req -newkey \
    	rsa:2048 \
    	-noenc \
        -keyout $$certname.key \
        -out $$certname.csr \
        -subj "/CN=root.kafkaexample.ru/OU=TEST/O=KafkaExample/L=Moscow/C=RU"

# convert the key to PKCS8, otherwise kafka/java cannot read it
.converttopkcs8:
	openssl pkcs8 \
  -topk8 \
  -in $$certname.key \
  -inform pem \
  -out $$certname-pkcs8.key \
  -outform pem \
  -nocrypt
	mv $$certname-pkcs8.key $$certname.key

# Sign the CSR with the root CA
SHELL = bash
.ONESHELL:
.signcsr:
	dns2=$(DNS2)
	openssl x509 -req \
      -CA root.crt \
      -CAkey root.key \
      -in $$dns1.csr \
      -out $$dns1-signed.crt \
      -sha256 \
      -days 365 \
      -CAcreateserial \
      -extensions v3_req \
      -extfile <(cat << EOF
	[v3_req]
	subjectAltName = @alt_names
	[alt_names]
	DNS.1 = $$dns1
	DNS.2 = $$dns2
	EOF
	)

# Combine private key and cert in one file
.combine:
	cat $$certname.key $$certname-signed.crt > $$certname-keypair.pem

# create properties for producer and consumer
SHELL = bash
.ONESHELL:
config:
	consumer_cert=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' consumer-signed.crt)
	consumer_key=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' consumer.key)
	truststore_cert=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' root.crt)
	DIR="$$( cd "$$( dirname "$${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
	echo
	echo "Create consumer.properties with inline certificate, private key and truststore"
	cat <<EOF > consumer.properties
	security.protocol=SSL
	ssl.keystore.type=PEM
	ssl.keystore.certificate.chain=$$consumer_cert
	ssl.keystore.key=$$consumer_key
	ssl.truststore.type=PEM
	ssl.truststore.certificates=$$truststore_cert
	EOF
	cat <<EOF > producer.properties
	security.protocol=SSL
	ssl.keystore.type=PEM
	ssl.keystore.location=$$DIR/producer-keypair.pem
	ssl.truststore.type=PEM
	ssl.truststore.location=$$DIR/root.crt
	EOF

SHELL = bash
.ONESHELL:
kafka3: # create kafka3 SSL Environment properties with inline certificate, private key and truststore
	kafka3_cert=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' kafka3-signed.crt)
	kafka3_key=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' kafka3.key)
	truststore_cert=$$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$$0;}' root.crt)
	echo
	echo "Create kafka3.properties with inline certificate, private key and truststore"
	cat <<EOF > kafka3.properties
	KAFKA_SSL_KEYSTORE_TYPE: PEM
	KAFKA_SSL_KEYSTORE_CERTIFICATE_CHAIN: $$kafka3_cert
	KAFKA_SSL_KEYSTORE_KEY: $$kafka3_key
	KAFKA_SSL_TRUESTORE_TYPE: PEM
	KAFKA_SSL_TRUESTORE_CERTIFICATES: $$truststore_cert
	EOF