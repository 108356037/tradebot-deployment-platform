package mq

import (
	"context"

	"encoding/json"
	"os"
	"time"

	"github.com/lithammer/shortuuid/v3"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
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

func PublishUpdateEvent(uuid string, targetResource TargetResourceTypes, resourceEventInfo map[string]string, updateInfo map[string]interface{}) {
	event := ResourceEvent{
		BasicEvent: BasicEvent{
			EventId:    shortuuid.New(),
			OccurredAt: time.Now().Format(time.RFC3339),
		},
		ResourceEventType:  ResourceUpdate,
		UserId:             uuid,
		TargetResourceType: targetResource,
		ResourceEventInfo:  resourceEventInfo,
		ResourceUpdateInfo: updateInfo,
	}
	mqPayload, _ := json.Marshal(event)
	err := PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)

	if err != nil {
		publishedToRetryQ := false
		for !publishedToRetryQ {
			err = PublishMsgRetryQ(mqPayload)
			if err == nil {
				log.Warnf("Published event %s to retry queue", event.EventId)
				publishedToRetryQ = true
			}
		}
	}

	log.Infof("Successfully published event %s", event.EventId)
}
