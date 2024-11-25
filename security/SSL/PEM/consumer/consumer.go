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

		"security.protocol":   "ssl",
		"ssl.key.pem":         "-----BEGIN PRIVATE KEY-----\n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\nsPEc8o372JlQJi0xIzJyTl8wRiCpWA==\n-----END PRIVATE KEY-----\n",
		"ssl.certificate.pem": "-----BEGIN CERTIFICATE-----\nMIIFqjCCA5KgAwIBAgIUQ+53zVdhCEe3C8D//gscj7tLr4QwDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjUwODM1MjFaFw0yNTExMjUwODM1MjFaMFcxETAPBgNV\nBAMMCGNvbnN1bWVyMQ0wCwYDVQQLDARURVNUMRUwEwYDVQQKDAxLYWZrYUV4YW1w\nbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQD1yFx4v6891YXe+KEuG5nNT9mqFZj6L2e4PBmnG7Z4\n0br1n4fDuyPkMHL5V0SCj7KaieAZ8NLYZmY2dcjBd1iXf7afh/7fBBhWlNBBeFJ6\n5IPG7kPjDfrKv8EFrxQ+Ssg5hNJZcznqQJzHu2jIfjEuzAbbQbvGoImBmu8wJDyB\nXpmdeSZphXdCRH2TdYh4lfH7cSfF6nv7ZBX3jRgH1dU9Z7KCp0nImHdvrTddAb/a\n5Ito5mXfwU/DWxMRPAJLeKi8UiLZx0Ee8MFXfhfyQXSQYzEf69oBiNTg6X39HGFg\nPBGD7T/wv6PEbl9m84PbCrqid+Lrcxddv7pzJDbjo3XuDDjwo6Y4ztQx0sJDp/Dy\nD/KLAPy9ktkIdoNvAB7+xO3f3KD5idzcniR87SG6CeQGKETvgPqhbH+LeUrAxcAT\nhdkranXMvb4UU34wqk2YZwM1vJn+OSNYO+O03pAZBLxorTk2Ezp2zJlKHfPXAEPl\nnPgFMVaO5A1swYybLylSpHMT5x0ClNO6TQfKLllgm5sdCU6Gx11c8oSf+V32y6fE\nH0WtHfb0FhzRPh91PJ34QWf4qjE5XgS5Ip82OP8kxasUPvmOH3sJK96reHMbn1I5\nfDDx4A87LK7dNFFSdshvJje6fFMyLNflG+xgIPFmId5TKj+KNXRJzYI1O7N8iXWN\nVQIDAQABo2IwYDAeBgNVHREEFzAVgghjb25zdW1lcoIJbG9jYWxob3N0MB0GA1Ud\nDgQWBBSRRJu32sbOjlB7yj3c5Cg+0EWx2DAfBgNVHSMEGDAWgBSfNqW2pOW/Bj1u\nurHLgtC59T8w1DANBgkqhkiG9w0BAQsFAAOCAgEAN1XrXt+91SQNXtULczlwxoNn\n8dCHnqgjzle8hntAZTfbYzaNmdIoVniOz/TZyH4//pElnVk+n+BVhXpbravRSege\nOUr3ZHkXlg8KwvtiRP2HErhWf3HXdYbT0izGCO0AKsx+QxB+wGA7PWzJkH7xFKws\nMAZZlZwWX5Owuhj+EtUy+V8hS9fh60nntPVTX1cZXiM0SofuC0Tidzuscucrn/4T\nrKTVvsZbgatPLJpdFbQiUBeVrWURpCLv7guEDEwZzgthcHKcpXHcJkvpA5Q1vLIT\nm+7H1eur1qMf053DnTyHPZ+wiU3XyszqembDzmRTtM1XscGpqu1iYAxgaO/eXJu9\nlCtx26yeVNWi30JObjjblA1pdMokafni+sTVr+UEHOvQe14i8MdEcDnxQ8lWpKqu\n6JCnMWF92tud9Kj6DH8Pp8jYa48AvKqlrLs8ICxWCZEiJc7MwNjIml8+Z6SWwJoO\nkU/lJZyU1LJqOYvzNZWsJP9mWEZZqecbi1pLiDuQqfHufYCfzOjG7wC8ZJjwUsGL\nzjgjsZqikGTwF/KOHUD0lq12tMMl4Pkhj+gqtwDIZUpFLrUesx7+YN5mPNgQLqYL\nsiL/QN24Y44JIv2L/keMpahFvIRDcL0u1UeIQpJJve+7mYaM+gHuIBCkZA7EixFv\n7JePkYXPAFiRk3UasLo=\n-----END CERTIFICATE-----\n",
		"ssl.ca.pem":          "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUNdD4SBgbTHkyfXnBTIRIYFfv6VgwDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjUwODM1MDhaFw0yNTExMjUwODM1MDhaMGMxHTAbBgNV\nBAMMFHJvb3Qua2Fma2FleGFtcGxlLnJ1MQ0wCwYDVQQLDARURVNUMRUwEwYDVQQK\nDAxLYWZrYUV4YW1wbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCzFj6qr1OnsitQUvbQImT065k2\nbtVbzmqZd23e1tYmAOO89QJ0+ltsoCxowUgFrixhzKZVQCiU8nLzoqtDlv5sZlEh\nUapEf1ol6XkRSTPze7G2OHfoqJNRByb7i4GE4HoYvJ1mW59s6S7ox2ZVay9DfJva\nbud2bCTmXO2R3tD8a/MY1CzTgtRBdkwXIA0N6RY9C320maFb2gmz5KoQzAm4MTIC\n2uwy/ibCArJ99P0dbULTN8ue0VnKQu/t38hE/xxl1LZLAmMQVACv05owT/UW+nWi\n1I1dkt2y7aYUhNuCdXlUV9HGOgieHNQNnSwXRYOrhmoc4gWr16igdpWkgOJgdM2p\nORxmJ+nf3ZCRet5g3UWYD9s5aFTjkjwqkzbxynCuEs2LwAUAsodxX3JzQLF6cqfj\nadKLmxGFYam7SzPwEy5YZ8ur+y9uoNy9YWi2spoKvzSSfZREWIelw5bAUvJsIggb\n6I3yG2sy4eyKlisM/Ug1QXSS4U2wauiH6vunMlSIMYR9g8Sl/xgl+/iYznsg2aK9\n9W6qiBcn9B4geokt4DscyjMRMYBjnEGzJGwpHCCKAWfp3dEbGuL3/jJPc+URN7xr\nVT0YHFxWrOXaGvHu0fiKAwVoiLvfV78PLsHUj4U5u+XetaUEFIaX0BlnWVBiI00e\n7BOnLNV1aAWH67pB8QIDAQABo1MwUTAdBgNVHQ4EFgQUnzaltqTlvwY9brqxy4LQ\nufU/MNQwHwYDVR0jBBgwFoAUnzaltqTlvwY9brqxy4LQufU/MNQwDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAesk1UuJQKvPgHgDRfAKAAMc3va1+\nrB6XFLbMp/EyaRxy+fXKLbEp542slMizaVWlnMGZVDJYyA2iuHPZY0Nb2P+5o9re\n986Au316HxpNqgbKkL3uUVQKSOYIfKkMIj+QWfnneCgRwSnYTGxM/iCMSinyjx70\ngC5jQ+XJQ8fMeRKme+blsk/2bqhyCfzqM+wGjZvl2z/vs9ky8jPL+YhpJJg873S2\nIiUm/Jhw2qYCZy0Z2nvLJD+lyrNcUTYOc6hw3NV7J3d4TeOSz6IeQnG4AmUwFbmd\nGBb38c1fVx62uxq0pmz7mgBTWGv3dLvXoIkICSW3+5+RHA5WfdFrB6Rhn6Qx7KBi\nNW1KV0ujfulRprBmjoaZfaJo4JXCLEgp5zl5Rt+kcaHuub7WbNdhI7jGY8DKIOp5\nO3TzyMtZmfmcAQBc5xS5nb64UauPbkU8Io2jmukTWq73nZsDlnSFk4Asahybrwp1\nAMYYq+5wUfMUuNw7adDmBEq3eCLqoZuNEYaZs7NVzndr+iy/7EUDxZ5+Mk7+SGdG\nzRRP1sbl/6IrF/4+MyJpvjdTNKt52H9RZ9I3YLtItjYg9tMsBXwQeWW75q3FRMKD\nU0/Pd/VxMbcKmZj0kVQcgQ5dRf7Xwlyv/igFYBNZyvK54rzv1qepT65LzwtEE2Em\nEsjnDbT1pM1oAGs=\n-----END CERTIFICATE-----\n",

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