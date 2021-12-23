package mq

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type EventTypes string

const (
	ResourceCreate EventTypes = "resource_create"
	ResourceDelete EventTypes = "resource_delete"
)

type BasicEvent struct {
	EventId    string     `json:"event_id"`
	EventType  EventTypes `json:"event_type"`
	OccurredAt string     `json:"occurred_at"`
}

type CreateEvent struct {
	BasicEvent     `json:",inline"`
	UserId         string `json:"user_id"`
	TargetResource string `json:"resource_type"`
}

type DeleteEvent struct {
	BasicEvent     `json:",inline"`
	UserId         string `json:"user_id"`
	TargetResource string `json:"resource_type"`
}

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
