package mq

import (
	"os"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

//type MsgHandler func(m kafka.Message) error

var (
	KafkaConsumer *kafka.Reader
	KafkaProducer *kafka.Writer
)

func InitSub() {
	maxWait, _ := strconv.Atoi(os.Getenv("MAX_WAIT"))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("BOOTSTRAP_SERVER")},
		GroupTopics: []string{os.Getenv("SUBSCRIBE_TOPIC")},
		GroupID:     os.Getenv("CONSUMER_GROUPID"),
		//Topic:    ,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  time.Second * time.Duration(maxWait),
	})

	r.SetOffset(kafka.LastOffset)
	KafkaConsumer = r
}

func InitPub() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVER")},
	})

	KafkaProducer = w
}
