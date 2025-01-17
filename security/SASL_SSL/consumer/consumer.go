package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	brokers               = "localhost:29092,localhost:29093,localhost:29094"
	topic                 = "topic"
	group                 = "MyGroup"
	defaultSessionTimeout = 6000
	ReadTimeout           = 1 * time.Second
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,

		//"security.protocol":   "ssl",
		"ssl.key.pem":         "-----BEGIN PRIVATE KEY-----\n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\nzsYvYq6WW16N81UXs4z7ikaH/oXEmDk=\n-----END PRIVATE KEY-----\n",
		"ssl.certificate.pem": "-----BEGIN CERTIFICATE-----\nMIIFqjCCA5KgAwIBAgIUF02l4XN/8vTiNT0wfhqvkpeeKoEwDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjcxMDA2NTJaFw0yNTExMjcxMDA2NTJaMFcxETAPBgNV\nBAMMCGNvbnN1bWVyMQ0wCwYDVQQLDARURVNUMRUwEwYDVQQKDAxLYWZrYUV4YW1w\nbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQCwXmxyabG+ykN4DrpNPJBEv5t2ady8BhlHtHVKNLK3\nHR6vL+E+8xN5ebgXEL3FIr2b+FSk4eV4J4wqJl8GYRf0zbuj9rZfTJE83/GMuu4R\nLWXfghK1IRQmKPxm9sGyiaROFVpBNPYTABB0UjkxI1PZB1U4jT6HsbE3nfjg3X34\nIdgk/3fwl8b2wao3ujj9GJOm883Sgv+wTHp9j2PnakYT+W7dBh8V3QgJR4Yq/J3L\nfnNZFuF0qmiw90IuLcXiViUtgPLVkIc46h0SgwM1ma0CEb0fkkGk4HrnHSUf+FQG\nbmCMm2cRXYm1lYCJ//4VwcfgGVg4s0DuhFrjIsSJO37c3mArPEsUWLnw/RnZCU3r\n5bDgjb03cJWBy4+o4yBantKPZl7vciILGDuXd9pkG+1WCrWL9+2ASeKgxuvrAgV7\nuiEU/k1Cuv/yWOtP0JpW6COzdyWow2ouUg2qOo0CUPav3PYYdHRvvXLWSFhUqOqF\nxsRE617SDbsK/YeajJEMpZtsm/N0qg1ucMSbnnwa2TUo9+gK3hes4zptHVdCHOcy\nyaHrH8HpkeCrD67oTWcBCO6O/lfyJYWD+c4kgdC3WW0PjhJiw/ANPck/6rr3kCr8\nYCIhP3XA9LqGmfNgG/dYGU9ZAGAxegLAehkksA6jWKVF0il2BnJHHE5QXJ5gtlBi\nFQIDAQABo2IwYDAeBgNVHREEFzAVgghjb25zdW1lcoIJbG9jYWxob3N0MB0GA1Ud\nDgQWBBRq3S0i+pDex5fQ/H/+kftN2GZ7STAfBgNVHSMEGDAWgBRZepXdNq1FWmRL\nLFtoPjLNIcc8TTANBgkqhkiG9w0BAQsFAAOCAgEAKeFZUuMrF8A8lR1j6l90O/sY\nBK/6UobcoSypSS5rOe6Xv5+sB0+PryjlGD3eEUDiIPbjYiMnzwDhPdGukPRTcLFu\nYGWXlZUWJWe6iUBjwGB+7BsEF2FjhJjcL7NZQg0wAkK/Cvfe1PhHTqxsh5q090wF\nG3qG1G3S5i0smUQaWe5wmc3Y88X9Ko58i3NKL7Sf+R8fJnX+7DzCWIBby/Cb+s9Z\nTypkCBVJ+xP8EZ5PcA5pXEYyKkSc8sFvtoFmg1tlYIurJJ6LjNe+8MxOzs9+ZWqN\ngwoLNP19M4KpTZi5uVXgjmijbuEED80mesOVwqCZQ9WPkVqkh1rfJasabWs8yb/S\nEkRpRTQ5EBfXEoZoegXxRx4C2Otb2ECfgaTFBkORiGdMQ5kC6TZPTtKPuTIPB92s\nG82+BMwjfLhhfv2In7uZlKEPLgHuAH3WevISBNoMNVPSylFS+2ZTR9XHf46cPmrN\nHsGvapqhPW+TR6cpkWE0dVikThlrEb3onzhdsS3kvlRCvX3MJpofuE1nc8TF8ggC\ndh8jW8zvKWQ41EL3kf/e0JjF0k6oppCWb+XvU42OBfUYYH+Ov65W2EPREOfBlEiW\nBAaxlCXLlOUte4MRD8XjoacZX3x7uqeyhwg61PEJA2+2S5FPXFmaGQkcUhYK0JS5\nimTMhvC/Kc8fPd9rWUg=\n-----END CERTIFICATE-----\n",
		"ssl.ca.pem":          "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUBa/ywVnqY/RjvrJUo+01FlMRR3wwDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjcxMDA2MzBaFw0yNTExMjcxMDA2MzBaMGMxHTAbBgNV\nBAMMFHJvb3Qua2Fma2FleGFtcGxlLnJ1MQ0wCwYDVQQLDARURVNUMRUwEwYDVQQK\nDAxLYWZrYUV4YW1wbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQD3HZxfKyzXyUhIYwFTuOeQSbrh\ndFdhw4n+o/DOStUtgZoR0EVGpBdfV++SSoNFHTsU4W81e1jp5hNyzH3N38Zlcu4j\n2Ci2mdK+jODO9vGrNWAvkrs/Th62JiADSqLWx07sTjutELFHUP39Ujha3qxrKaDG\nZwqvo7jevP8UkkYnTauJ/8eRfYMD5zuR6EM5A/Pj7EGcGCxByvE9OP5Puq1fqJ/l\nR8gwBT92XrqlfhYB5C7p643/lETtbN048eI0M6C5daO8BbhCzG5T30PFx0/Gwxzh\ntMse1RzYdzvpgl+hl9svvx0OLj7CEe0AvL7IYAitfAk18md1YiqUcnvXcS6zJScj\npOWjwryj9aiIx1PCjejSrQUrr1tYkpCZsWCYG1f66+se18h2K1ZoKDxVJTHab8Ju\niAVswOZoE4pUnuOQkr2oa2lvbWqc0BDk+uMCVLsVeEUW3/MALKO1ScDzQwhVSRFf\ny72rY79gmnaqRqJe/TaM5OTOn6qvieg8z4wEgZ0Q0wzjCr54w85eHGbsL8j9K1KF\n0qInE45kzvBe5XkDK1PDpizkO2LH0LiDmWyNrxUE+WgmoDHaJLNY8n7clvSN9p5B\nbssgpxmihAyr3K88A7E7xSk0git6wjnlNZ83PxURkHDy6CNmFDqV6Q+S0x53Fep2\nYkmqdZXK1rdRhy8MXwIDAQABo1MwUTAdBgNVHQ4EFgQUWXqV3TatRVpkSyxbaD4y\nzSHHPE0wHwYDVR0jBBgwFoAUWXqV3TatRVpkSyxbaD4yzSHHPE0wDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEA6NT4di5g6rqhdnPtgROEioi98aE/\ns6eP0qRl4c/dSG7G/0uDUICCQ0VH5TG0LYeImrHKR534TvYD+BqocvJdZnpM6LHC\nShk29hHhvQ0n9T2FZY/SIm2DtNkHcdNBQLlLlq0zVPsKY/IMyfoOqa1C1MSJAc7O\n2f8HVtG8Gxmi+FDqkz/Ye4KntB5NIMad6MTnG+ZB9bBTnjlkkDI6Mmzz5C9wwMs7\n8YX3NpEJGAu804bF1woO+vhFpLCnW/8VpnEz+/2jj/8zTgxUAjqu1n2QJE5jxYyl\nhB33bHxfASdCsb+XxLGb5nU53pr/+c5rB69MraISfnfaq4TjXujxRosbcDVCkZoW\nNI+ImSxiS9YwCQcCRrtkFIQ/zUFhJn6t1y7JU87bFjSxYKH+QTWTqqR6JwS+ORm8\nnyjY+0cJa3BZB8Fvmq9V6A22D71c53XpNrB2T36jarsYjGBYlWfjPFya0ZM2u8na\npw6fdC3Err4YFkZUArHptzw09Rbp6yZWBN8jaXq1BtQOb93hDkhoC5tkktbH+R+E\nLdobLqc1LLX8f7ctYTaWFqITiYXaNwrRPoEAXMfIyitHhT5JMCm84r9Sdbo2XLZE\nKPEuWvooJ3Wt0H2cAs5llnxsFxaSz3hefio+6Tjam7Ijnd95RKwU2qcScFCK85Wu\nIaZS2pVXBFFzSOg=\n-----END CERTIFICATE-----\n",

		"security.protocol": "sasl_ssl",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "consumer",
		"sasl.password":     "consumer-secret",

		"group.id":                group,
		"session.timeout.ms":      defaultSessionTimeout,
		"enable.auto.commit":      "true",
		"auto.commit.interval.ms": 1000,
		"auto.offset.reset":       "earliest", // earliest - сообщения с commit offset; latest - новые сообщение
	})
	if err != nil {
		slog.Error("Failed to create consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		slog.Error("Failed to subscribe to topic: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-stop:
			break loop
		default:
			msg, err := consumer.ReadMessage(ReadTimeout) // read message with timeout
			if err == nil {
				fmt.Printf("Message on topic %s: key = %s, value = %s, offset = %d\n", *msg.TopicPartition.Topic, string(msg.Key), string(msg.Value), msg.TopicPartition.Offset)
			} else if !err.(kafka.Error).IsTimeout() { // TODO process timeout
				slog.Error("Consumer error: ", slog.String("error", err.Error()))
			}
		}
	}
}
