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
		"ssl.key.pem":         "-----BEGIN PRIVATE KEY-----\n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\nAxxPNAfMU/e+80561fSjkPbCfR+XwT/4\n-----END PRIVATE KEY-----\n",
		"ssl.certificate.pem": "-----BEGIN CERTIFICATE-----\nMIIFqjCCA5KgAwIBAgIUI0vz9PWya77qW2Ak/c9e98KqVZowDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjUxNzA3MzdaFw0yNTExMjUxNzA3MzdaMFcxETAPBgNV\nBAMMCGNvbnN1bWVyMQ0wCwYDVQQLDARURVNUMRUwEwYDVQQKDAxLYWZrYUV4YW1w\nbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQDK1+bz/IVCRrwEaYtlz487+PkLGFnJM0W973bW4Hi7\nM0vREXQtZ5L1anU8R9oqhVqoGrexqSc8j/NXyFgVxGC4bzspdIgZ4hG+phWBrk3S\n01Kp3tcQJnYXJ2K4mrKZykNsJjaD/JPm0Sp4hH4OK8Hzs5Alh68u4axTrJ1suIPg\n7Q4tBJuu3BMLtrfbcuMISTvhacoFwBFuRp3SQmZ6VOF+G0z/Lg6gAL8RZd547rcI\nhdtynbWNiPETrFfcMa9bXeLI8fdVIqwGZArsvtMAHJOPkCWPcAm173ECjUoJ7P4E\nexkOphqt4bSARYZQ3JtqJXPv78tvKYMIl9rlJ961TjE1sTwO+C9m1UMS8cfJhWxg\n1CP1tepfGErHyRJ7wqMcm277pwRE5GLsg628ABOCFynIYWOrrF3UYMBqgS3+L5/6\nGwOUrVpaok6oQHgtwN5J0Xagpu2G0mBvZJERbiKwRjo6THGJcBNsudvDlTKL+RtW\nDA78DrA+e44aZjZO1BFYmvmbKhQzG44H433gKTWTQ/HrU1bSFj6RKZZ9wsoiu6w1\n8MgX3mpq9eXF+3dnk6bI/N/qVTExS8H12eLH/u3D8EjETFK/gxGcxj0RScvsF9bb\nTdPOoNAKMddNNKOe1t9vNQFo3tbXNywTqGRypanx8iSnpryP+jWYRAJS/sv6r8jn\nTwIDAQABo2IwYDAeBgNVHREEFzAVgghjb25zdW1lcoIJbG9jYWxob3N0MB0GA1Ud\nDgQWBBQUTXPm3qc7ArDtQWju/RqJkeXW/zAfBgNVHSMEGDAWgBTBDY5IX1EOqLdU\nEggTNolnEi3q4zANBgkqhkiG9w0BAQsFAAOCAgEANxZqDXTGcMyC8oEzZx6lsvOk\nieP5SjC1HeuRY1EwYdEuHDzKPoNG0661LrLQf/Qe4MJg9dbbGSdk4lxQvGvYrcqa\nHTih27FqpTodFogOda3Tv7OldzrLYwnNvYLa8JYTLVR+dINvcwAjaeiawynP/ZAJ\nYCeoFNs+XhITbEPOAlKsD/P5WpBSpvAF2/qhe6evWB9sMcjZ+Cb8vLi8IT22YNFe\ntXP1Rq4YPKeoiCquwF1DK6Yku/aQGuDWf9WUWqlZZ8V0vgf6DNm4gYF+FwT2J3xK\n4g5P/ZEj5YwP9pLf+qD5LEPKWBUP5MtzLSaxWcuK1Hf/XTf/RN95ok4IlP+WgtKi\nVTTtYDm0g10fdcpZcjDNCNOmD1hW1nlnp0Wr0o6QV3uLtft31mSy/IUmJbxm+9/h\noP9Hl11d+TUWVV+XQXKv++QSjVM678XyV6f/wcUCd/RkXBc9PI6JnYms2sZRClyY\n7h0AQs+V0Uo9QFx7EgXFFPH55FAOfLILzwsKTzD3kVDlGRuRmJl42BhsygyD0Rmo\n/y7G0tqc/eiPOPt7UNuLxYiwvaFyMBwjke2RkLwjMP/j7iZ4maI7Gik6aNxp9/89\noC19Rh62r7nSGjfMc8xNVkYilz4gSmyP9ajb4yfCpUNRaxiRLo7lgMaGICUh1FIv\nXoCdE5mkrx0SzwnQkPs=\n-----END CERTIFICATE-----\n",
		"ssl.ca.pem":          "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUGSF9DTyk3rih/9zVdLZnsRHg4sAwDQYJKoZIhvcNAQEL\nBQAwYzEdMBsGA1UEAwwUcm9vdC5rYWZrYWV4YW1wbGUucnUxDTALBgNVBAsMBFRF\nU1QxFTATBgNVBAoMDEthZmthRXhhbXBsZTEPMA0GA1UEBwwGTW9zY293MQswCQYD\nVQQGEwJSVTAeFw0yNDExMjUxNzA3MjBaFw0yNTExMjUxNzA3MjBaMGMxHTAbBgNV\nBAMMFHJvb3Qua2Fma2FleGFtcGxlLnJ1MQ0wCwYDVQQLDARURVNUMRUwEwYDVQQK\nDAxLYWZrYUV4YW1wbGUxDzANBgNVBAcMBk1vc2NvdzELMAkGA1UEBhMCUlUwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCQXKHDLJ4vBb60g3iJSWtk/4G5\nV+bF23QyV/Xb0ZQUxcVuTI0gKY+p3crMrw44KEb6NpSHLDzpogjUopTnjY+7Vwl1\nz8/C89TT5H/X2Hfl1ekso/Cfa6ZHUTwdZGbcMw9rqX2pCdQT4FiiVGZLYXWkVZyR\nYjjImYe4SlWqdk/dhtQ3rppTlJjY2Y8ONCuXwbUraIAgdM4WfNWbEm1TKvKzKett\nDUcWWBzqbHafAu73EMkYRfSl88NnxQILgxPBU7lwQy4tIPbiDChBUtIymuZX84m0\nCCH/o4wsbV3jLYtggn6XUhx9ndHMvn4fV1ctmczEc1TfiYGatYiZbz7W6so+/vDV\njGnM4CTMZRsXj04FUIaNq6PhxAM63346t9zdEesW3+lMfYZt8CzdDiSvNk7pcdfJ\nFkJSZTaiyCl9RF7+XiaW+sqMaZw0qUKDNwdqPTl3GKNLoPm47JiZWhMOrKbXsu70\nlIEErvaZedGsysTUQa06qQdcajuwmzzb++Q/lkyk6i7wFr88C45dMSUbMCrEfZqh\nfStW28GMRImuT1XcQ5/V5UNQnz3F76Y0/0UBqOQqSGG0hw2CrCZv9qT+XuhOx4tb\nK5th15KRxuTapXUKaOkMFfo8JoMTFINKMPlltXKOcr2hW/mJzzSsz3xIvrC0CB9I\ntiaT8rO/e05HDDL9qQIDAQABo1MwUTAdBgNVHQ4EFgQUwQ2OSF9RDqi3VBIIEzaJ\nZxIt6uMwHwYDVR0jBBgwFoAUwQ2OSF9RDqi3VBIIEzaJZxIt6uMwDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAd3Z6rW41m3m9xXeVup4tOodAng3X\nxNkzzj+Wa4qyBK0aBuEDYHW7JIRaMTdTeDYo15s7nK0oSWkw1lEFguxvYjGxkUzq\n5SrcRMe8+hxf3tORgia7rgbSkU5neY7jt72koLtYDVTJQgElagH9/VtbZdryn+MA\nmit+YvVnqI2FGRk/0395TTKt5A076ircQSjHxgrTsNfsBJKJElwQxGGiXDp63QMP\nduy3KNNV6fDPexAkmm9Gs7yl5AACtFRsJarxYMz8cmV4fBYbOGN231Bx8TIu82ym\n6NDrUbdaeE6sKfYSGo9kAcSMPTtMdgnAtrNGDYmQM68FAvrxaKPUAswnKS9ERg9o\np3IvVGvCY6wtiwM++n71xrkCkCF26mOSJ58mxSEdxysPd8GKDZ7sfQf14geVOL9c\ngWoCtpOYFW9wj1zJV56ox8sWFW+tiF4haIXgHkvWfZgId72kAfaK7B2AYTMGo4Yj\nNs7F/haCMelJvYep1GIlKiUvcYotQW99qep48sCGxd2SsjyNfHihdMCpCrsELBce\nUhFOHCZkGXzy0aPZmE+YFni/A9KDnx7ZlY2mJMnXvW+6UTiUo9nAHgcoCFXmFUKb\n0BF9KXVdEfNfYnqY4oRBQOu+AkvnqNq0Dc070MyqZlNSkE/e0Al2rzaqKs1sHVrf\nmzUQ0CFt1XWBSaE=\n-----END CERTIFICATE-----\n",

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
