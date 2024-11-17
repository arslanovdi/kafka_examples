package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-jose/go-jose/v4/json"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var brokers = "127.0.0.1:29092,127.0.0.1:29093,127.0.0.1:29094"

const topic = "stats"

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":      brokers,
		"group.id":               "group",
		"session.timeout.ms":     6000,
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 5000,
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
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Stats:
				var stats1 map[string]interface{} // Example with map[string]interface{}
				json.Unmarshal([]byte(e.String()), &stats1)
				fmt.Println("name:", stats1["name"])
				fmt.Println("client_id:", stats1["client_id"])

				stats2, err := Unmarshal([]byte(e.String())) // Example with unmarshal to struct
				if err != nil {
					slog.Error("Unmasrhal error", slog.String("error", err.Error()))
				}
				fmt.Println(stats2.Brokers.Fields[0].Name)
				fmt.Println(stats2.Brokers.Fields[1].Name)
				fmt.Println(stats2.Brokers.Fields[2].Name)
			}
		}
	}
}

func Unmarshal(data []byte) (*Stats, error) {
	/*
		Имя брокеров разное в зависимости от проекта.
		в JSON stats один из ключей поля - имя брокера, для анмаршаллинга вырезаем эти ключи. Собираем параметры брокеров в слайс под ключом "fields".
	*/

	str := string(data)
	a := strings.Index(str, "brokers\":{")
	a += 11 // "brokers\"{ :"

	b := strings.Index(str[a:], "{ \"name\":") // первый брокер

	str = str[:a] + "\"fields\": [ " + str[a+b:] // заключаем брокеров в фигурные скобки, чтобы анмаршалить в слайс
	a += 13
	igroup := strings.Index(str[a:], "\"GroupCoordinator\"") // конец блока брокеров

	b = strings.Index(str[a:], "{ \"name\":")
	for b < igroup && b != -1 { // Ищем всех брокеров
		i := b
		for str[i+a] != ',' {
			i--
		}
		i += 2 // ", "
		str = str[:i+a] + str[b+a:]
		a = i + a + 1
		igroup = strings.Index(str[a:], "\"GroupCoordinator\"") // конец блока брокеров
		b = strings.Index(str[a:], "{ \"name\":")
	}

	str = str[:a+igroup-2] + "] " + str[a+igroup-2:]

	var stats Stats

	// Анмаршал исправленного JSON
	err := json.Unmarshal([]byte(str), &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// Stats Пример структуры. НеПолное описание тут: https://github.com/confluentinc/librdkafka/blob/master/STATISTICS.md
type Stats struct {
	Name             string `json:"name"`
	ClientId         string `json:"client_id"`
	Type             string `json:"type"`
	Ts               int64  `json:"ts"`
	Time             int64  `json:"time"`
	Age              int64  `json:"age"`
	Replyq           uint64 `json:"replyq"`
	MsgCnt           uint64 `json:"msg_cnt"`
	MsgSize          uint64 `json:"msg_size"`
	MsgMax           int64  `json:"msg_max"`
	MsgSizeMax       int64  `json:"msg_size_max"`
	Tx               int64  `json:"tx"`
	TxBytes          int64  `json:"tx_bytes"`
	Rx               int64  `json:"rx"`
	RxBytes          int64  `json:"rx_bytes"`
	Txmsgs           int64  `json:"txmsgs"`
	TxmsgBytes       int64  `json:"txmsg_bytes"`
	Rxmsgs           int64  `json:"rxmsgs"`
	RxmsgBytes       int64  `json:"rxmsg_bytes"`
	SimpleCnt        int64  `json:"simple_cnt"`
	MetadataCacheCnt uint64 `json:"metadata_cache_cnt"`
	Brokers          struct {
		Fields           []Field `json:"fields"`
		GroupCoordinator struct {
			Name           string `json:"name" json:"name,omitempty"`
			Nodeid         int64  `json:"nodeid" json:"nodeid,omitempty"`
			Nodename       string `json:"nodename" json:"nodename,omitempty"`
			Source         string `json:"source" json:"source,omitempty"`
			State          string `json:"state" json:"state,omitempty"`
			Stateage       uint64 `json:"stateage" json:"stateage,omitempty"`
			OutbufCnt      uint64 `json:"outbuf_cnt" json:"outbuf_cnt,omitempty"`
			OutbufMsgCnt   uint64 `json:"outbuf_msg_cnt" json:"outbuf_msg_cnt,omitempty"`
			WaitrespCnt    uint64 `json:"waitresp_cnt" json:"waitresp_cnt,omitempty"`
			WaitrespMsgCnt uint64 `json:"waitresp_msg_cnt" json:"waitresp_msg_cnt,omitempty"`
			Tx             int64  `json:"tx" json:"tx,omitempty"`
			Txbytes        int64  `json:"txbytes" json:"txbytes,omitempty"`
			Txerrs         int64  `json:"txerrs" json:"txerrs,omitempty"`
			Txretries      int64  `json:"txretries" json:"txretries,omitempty"`
			Txidle         int64  `json:"txidle" json:"txidle,omitempty"`
			ReqTimeouts    int64  `json:"req_timeouts" json:"req_timeouts,omitempty"`
			Rx             int64  `json:"rx" json:"rx,omitempty"`
			Rxbytes        int64  `json:"rxbytes" json:"rxbytes,omitempty"`
			Rxerrs         int64  `json:"rxerrs" json:"rxerrs,omitempty"`
			Rxcorriderrs   int64  `json:"rxcorriderrs" json:"rxcorriderrs,omitempty"`
			Rxpartial      int64  `json:"rxpartial" json:"rxpartial,omitempty"`
			Rxidle         int64  `json:"rxidle" json:"rxidle,omitempty"`
			Req            struct {
				Fetch                     int `json:"Fetch" json:"fetch,omitempty"`
				ListOffsets               int `json:"ListOffsets" json:"list_offsets,omitempty"`
				Metadata                  int `json:"Metadata" json:"metadata,omitempty"`
				OffsetCommit              int `json:"OffsetCommit" json:"offset_commit,omitempty"`
				OffsetFetch               int `json:"OffsetFetch" json:"offset_fetch,omitempty"`
				FindCoordinator           int `json:"FindCoordinator" json:"find_coordinator,omitempty"`
				JoinGroup                 int `json:"JoinGroup" json:"join_group,omitempty"`
				Heartbeat                 int `json:"Heartbeat" json:"heartbeat,omitempty"`
				LeaveGroup                int `json:"LeaveGroup" json:"leave_group,omitempty"`
				SyncGroup                 int `json:"SyncGroup" json:"sync_group,omitempty"`
				SaslHandshake             int `json:"SaslHandshake" json:"sasl_handshake,omitempty"`
				ApiVersion                int `json:"ApiVersion" json:"api_version,omitempty"`
				SaslAuthenticate          int `json:"SaslAuthenticate" json:"sasl_authenticate,omitempty"`
				DescribeCluster           int `json:"DescribeCluster" json:"describe_cluster,omitempty"`
				DescribeProducers         int `json:"DescribeProducers" json:"describe_producers,omitempty"`
				Unknown62                 int `json:"Unknown-62?" json:"unknown_62,omitempty"`
				DescribeTransactions      int `json:"DescribeTransactions" json:"describe_transactions,omitempty"`
				ListTransactions          int `json:"ListTransactions" json:"list_transactions,omitempty"`
				Unknown69                 int `json:"Unknown-69?" json:"unknown_69,omitempty"`
				Unknown70                 int `json:"Unknown-70?" json:"unknown_70,omitempty"`
				GetTelemetrySubscriptions int `json:"GetTelemetrySubscriptions" json:"get_telemetry_subscriptions,omitempty"`
				PushTelemetry             int `json:"PushTelemetry" json:"push_telemetry,omitempty"`
				Unknown73                 int `json:"Unknown-73?" json:"unknown_73,omitempty"`
			} `json:"req" json:"req"`

			ZbufGrow    int64 `json:"zbuf_grow" json:"zbuf_grow,omitempty"`
			BufGrow     int64 `json:"buf_grow" json:"buf_grow,omitempty"`
			Wakeups     int64 `json:"wakeups" json:"wakeups,omitempty"`
			Connects    int64 `json:"connects" json:"connects,omitempty"`
			Disconnects int64 `json:"disconnects" json:"disconnects,omitempty"`
			IntLatency  struct {
				Min        int `json:"min" json:"min,omitempty"`
				Max        int `json:"max" json:"max,omitempty"`
				Avg        int `json:"avg" json:"avg,omitempty"`
				Sum        int `json:"sum" json:"sum,omitempty"`
				Stddev     int `json:"stddev" json:"stddev,omitempty"`
				P50        int `json:"p50" json:"p_50,omitempty"`
				P75        int `json:"p75" json:"p_75,omitempty"`
				P90        int `json:"p90" json:"p_90,omitempty"`
				P95        int `json:"p95" json:"p_95,omitempty"`
				P99        int `json:"p99" json:"p_99,omitempty"`
				P9999      int `json:"p99_99" json:"p_9999,omitempty"`
				Outofrange int `json:"outofrange" json:"outofrange,omitempty"`
				Hdrsize    int `json:"hdrsize" json:"hdrsize,omitempty"`
				Cnt        int `json:"cnt" json:"cnt,omitempty"`
			} `json:"int_latency" json:"int_latency"`
			OutbufLatency struct {
				Min        int `json:"min" json:"min,omitempty"`
				Max        int `json:"max" json:"max,omitempty"`
				Avg        int `json:"avg" json:"avg,omitempty"`
				Sum        int `json:"sum" json:"sum,omitempty"`
				Stddev     int `json:"stddev" json:"stddev,omitempty"`
				P50        int `json:"p50" json:"p_50,omitempty"`
				P75        int `json:"p75" json:"p_75,omitempty"`
				P90        int `json:"p90" json:"p_90,omitempty"`
				P95        int `json:"p95" json:"p_95,omitempty"`
				P99        int `json:"p99" json:"p_99,omitempty"`
				P9999      int `json:"p99_99" json:"p_9999,omitempty"`
				Outofrange int `json:"outofrange" json:"outofrange,omitempty"`
				Hdrsize    int `json:"hdrsize" json:"hdrsize,omitempty"`
				Cnt        int `json:"cnt" json:"cnt,omitempty"`
			} `json:"outbuf_latency" json:"outbuf_latency"`
			Rtt struct {
				Min        int `json:"min" json:"min,omitempty"`
				Max        int `json:"max" json:"max,omitempty"`
				Avg        int `json:"avg" json:"avg,omitempty"`
				Sum        int `json:"sum" json:"sum,omitempty"`
				Stddev     int `json:"stddev" json:"stddev,omitempty"`
				P50        int `json:"p50" json:"p_50,omitempty"`
				P75        int `json:"p75" json:"p_75,omitempty"`
				P90        int `json:"p90" json:"p_90,omitempty"`
				P95        int `json:"p95" json:"p_95,omitempty"`
				P99        int `json:"p99" json:"p_99,omitempty"`
				P9999      int `json:"p99_99" json:"p_9999,omitempty"`
				Outofrange int `json:"outofrange" json:"outofrange,omitempty"`
				Hdrsize    int `json:"hdrsize" json:"hdrsize,omitempty"`
				Cnt        int `json:"cnt" json:"cnt,omitempty"`
			} `json:"rtt" json:"rtt"`
			Throttle struct {
				Min        int `json:"min" json:"min,omitempty"`
				Max        int `json:"max" json:"max,omitempty"`
				Avg        int `json:"avg" json:"avg,omitempty"`
				Sum        int `json:"sum" json:"sum,omitempty"`
				Stddev     int `json:"stddev" json:"stddev,omitempty"`
				P50        int `json:"p50" json:"p_50,omitempty"`
				P75        int `json:"p75" json:"p_75,omitempty"`
				P90        int `json:"p90" json:"p_90,omitempty"`
				P95        int `json:"p95" json:"p_95,omitempty"`
				P99        int `json:"p99" json:"p_99,omitempty"`
				P9999      int `json:"p99_99" json:"p_9999,omitempty"`
				Outofrange int `json:"outofrange" json:"outofrange,omitempty"`
				Hdrsize    int `json:"hdrsize" json:"hdrsize,omitempty"`
				Cnt        int `json:"cnt" json:"cnt,omitempty"`
			} `json:"throttle" json:"throttle"`

			Toppars struct {
			} `json:"toppars" json:"toppars"`
		} `json:"GroupCoordinator" json:"group_coordinator"`
	} `json:"brokers"`
	Topics struct {
		Topic       string   `json:"topic,omitempty"`
		Age         uint64   `json:"age,omitempty"`
		MetadataAge uint64   `json:"metadata_age,omitempty"`
		Batchsize   struct{} `json:"batchsize"`
		Batchcnt    struct{} `json:"batchcnt"`
		Partitions  struct{} `json:"partitions"`
	} `json:"topics"`
	Cgrp struct {
		State           string `json:"state"`
		Stateage        int    `json:"stateage"`
		JoinState       string `json:"join_state"`
		RebalanceAge    int    `json:"rebalance_age"`
		RebalanceCnt    int    `json:"rebalance_cnt"`
		RebalanceReason string `json:"rebalance_reason"`
		AssignmentSize  int    `json:"assignment_size"`
	} `json:"cgrp"`
}

type Field struct {
	Name           string `json:"name"`
	Nodeid         int    `json:"nodeid"`
	Nodename       string `json:"nodename"`
	Source         string `json:"source"`
	State          string `json:"state"`
	Stateage       int    `json:"stateage"`
	OutbufCnt      int    `json:"outbuf_cnt"`
	OutbufMsgCnt   int    `json:"outbuf_msg_cnt"`
	WaitrespCnt    int    `json:"waitresp_cnt"`
	WaitrespMsgCnt int    `json:"waitresp_msg_cnt"`
	Tx             int    `json:"tx"`
	Txbytes        int    `json:"txbytes"`
	Txerrs         int    `json:"txerrs"`
	Txretries      int    `json:"txretries"`
	Txidle         int    `json:"txidle"`
	ReqTimeouts    int    `json:"req_timeouts"`
	Rx             int    `json:"rx"`
	Rxbytes        int    `json:"rxbytes"`
	Rxerrs         int    `json:"rxerrs"`
	Rxcorriderrs   int    `json:"rxcorriderrs"`
	Rxpartial      int    `json:"rxpartial"`
	Rxidle         int    `json:"rxidle"`
	ZbufGrow       int    `json:"zbuf_grow"`
	BufGrow        int    `json:"buf_grow"`
	Wakeups        int    `json:"wakeups"`
	Connects       int    `json:"connects"`
	Disconnects    int    `json:"disconnects"`
	IntLatency     struct {
		Min        int `json:"min"`
		Max        int `json:"max"`
		Avg        int `json:"avg"`
		Sum        int `json:"sum"`
		Stddev     int `json:"stddev"`
		P50        int `json:"p50"`
		P75        int `json:"p75"`
		P90        int `json:"p90"`
		P95        int `json:"p95"`
		P99        int `json:"p99"`
		P9999      int `json:"p99_99"`
		Outofrange int `json:"outofrange"`
		Hdrsize    int `json:"hdrsize"`
		Cnt        int `json:"cnt"`
	} `json:"int_latency"`
	OutbufLatency struct {
		Min        int `json:"min"`
		Max        int `json:"max"`
		Avg        int `json:"avg"`
		Sum        int `json:"sum"`
		Stddev     int `json:"stddev"`
		P50        int `json:"p50"`
		P75        int `json:"p75"`
		P90        int `json:"p90"`
		P95        int `json:"p95"`
		P99        int `json:"p99"`
		P9999      int `json:"p99_99"`
		Outofrange int `json:"outofrange"`
		Hdrsize    int `json:"hdrsize"`
		Cnt        int `json:"cnt"`
	} `json:"outbuf_latency"`
	Rtt struct {
		Min        int `json:"min"`
		Max        int `json:"max"`
		Avg        int `json:"avg"`
		Sum        int `json:"sum"`
		Stddev     int `json:"stddev"`
		P50        int `json:"p50"`
		P75        int `json:"p75"`
		P90        int `json:"p90"`
		P95        int `json:"p95"`
		P99        int `json:"p99"`
		P9999      int `json:"p99_99"`
		Outofrange int `json:"outofrange"`
		Hdrsize    int `json:"hdrsize"`
		Cnt        int `json:"cnt"`
	} `json:"rtt"`
	Throttle struct {
		Min        int `json:"min"`
		Max        int `json:"max"`
		Avg        int `json:"avg"`
		Sum        int `json:"sum"`
		Stddev     int `json:"stddev"`
		P50        int `json:"p50"`
		P75        int `json:"p75"`
		P90        int `json:"p90"`
		P95        int `json:"p95"`
		P99        int `json:"p99"`
		P9999      int `json:"p99_99"`
		Outofrange int `json:"outofrange"`
		Hdrsize    int `json:"hdrsize"`
		Cnt        int `json:"cnt"`
	} `json:"throttle"`
	Req struct {
		Fetch                     int `json:"Fetch"`
		ListOffsets               int `json:"ListOffsets"`
		Metadata                  int `json:"Metadata"`
		OffsetCommit              int `json:"OffsetCommit"`
		OffsetFetch               int `json:"OffsetFetch"`
		FindCoordinator           int `json:"FindCoordinator"`
		JoinGroup                 int `json:"JoinGroup"`
		Heartbeat                 int `json:"Heartbeat"`
		LeaveGroup                int `json:"LeaveGroup"`
		SyncGroup                 int `json:"SyncGroup"`
		SaslHandshake             int `json:"SaslHandshake"`
		ApiVersion                int `json:"ApiVersion"`
		SaslAuthenticate          int `json:"SaslAuthenticate"`
		DescribeCluster           int `json:"DescribeCluster"`
		DescribeProducers         int `json:"DescribeProducers"`
		Unknown62                 int `json:"Unknown-62?"`
		DescribeTransactions      int `json:"DescribeTransactions"`
		ListTransactions          int `json:"ListTransactions"`
		Unknown69                 int `json:"Unknown-69?"`
		Unknown70                 int `json:"Unknown-70?"`
		GetTelemetrySubscriptions int `json:"GetTelemetrySubscriptions"`
		PushTelemetry             int `json:"PushTelemetry"`
		Unknown73                 int `json:"Unknown-73?"`
	} `json:"req"`
	Toppars struct {
	} `json:"toppars"`
}
