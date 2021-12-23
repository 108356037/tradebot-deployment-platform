package mq

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

func PublishMsgNoKey(topic string, payload []byte) error {
	err := KafkaProducer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: payload,
	})
	if err != nil {
		return err
	}
	return nil
}

func PublishMsgRetryQ(payload []byte) error {
	err := KafkaProducer.WriteMessages(context.Background(), kafka.Message{
		Topic: "retry_queue",
		Value: payload,
	})
	if err != nil {
		return err
	}
	return nil
}
