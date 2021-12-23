package mq

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

type msgHandler func(m kafka.Message) error

func SimpleLogHandler(m kafka.Message) error {
	log.Printf("message at offset %d: %v = %v\n", m.Offset, string(m.Key), string(m.Value))
	return nil
}

func RunSubscribe(handlefunc msgHandler) {
	for {
		m, err := KafkaConsumer.FetchMessage(context.Background())
		if err != nil {
			break
		}
		if err := handlefunc(m); err != nil {
			break
		}
		KafkaConsumer.CommitMessages(context.Background(), m)
	}
}
