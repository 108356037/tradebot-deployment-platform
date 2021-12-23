// this file init the publisher, subscriber settings
package mq

import (
	"os"

	kafka "github.com/segmentio/kafka-go"
)

var (
	KafkaConsumerList []*kafka.Reader
	KafkaProducer     *kafka.Writer
)

func Init() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVER")},
	})

	KafkaProducer = w
}
