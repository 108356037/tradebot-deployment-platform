package mq

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/108356037/v1/strategy-manager/internal/models"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type kafkaMsgHandler func(m kafka.Message) error

func RegisterSubscriber(kafkaTopic string, handleFunc kafkaMsgHandler) {

	maxWait, _ := strconv.Atoi(os.Getenv("MAX_WAIT"))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("BOOTSTRAP_SERVER")},
		GroupTopics: []string{kafkaTopic}, // current topics: resource
		GroupID:     os.Getenv("CONSUMER_GROUPID"),
		//Topic:    ,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  time.Second * time.Duration(maxWait),
	})
	r.SetOffset(kafka.LastOffset)
	KafkaConsumerList = append(KafkaConsumerList, r)

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		if err := handleFunc(m); err != nil {
			log.Warnf("Handler func failed: %s", err.Error())
		}
		r.CommitMessages(context.Background(), m)
	}

}

func SimpleLogHandler(m kafka.Message) error {
	log.Printf("message at offset %d: %v = %v\n", m.Offset, string(m.Key), string(m.Value))
	return nil
}

func ResourceEventHandler(m kafka.Message) error {

	var resourceEvent ResourceEvent
	err := json.Unmarshal(m.Value, &resourceEvent)
	if err != nil {
		return err
	}

	switch resourceEvent.ResourceEventType {

	case ResourceCreate:

		if resourceEvent.TargetResourceType == Strategy {
			err = models.CreateStrategy(resourceEvent.UserId, resourceEvent.ResourceEventInfo["strategy"])
		} else if resourceEvent.TargetResourceType == Namespace {
			log.Infof("namespace %s created", resourceEvent.UserId)
		} else {
			log.Infof("Nothing doned for event %v", resourceEvent)
		}

	case ResourceDelete:
		if resourceEvent.TargetResourceType == Strategy {
			_, err = models.DeleteUserStrategy(resourceEvent.UserId, resourceEvent.ResourceEventInfo["strategy"])
		} else if resourceEvent.TargetResourceType == Namespace {
			err = models.DeleteUserStrategies(resourceEvent.UserId)
		} else {
			log.Infof("Nothing doned for event %v", resourceEvent)
		}

	case ResourceUpdate:
		if resourceEvent.TargetResourceType == Strategy {
			_, err = models.UpdateUserStrategy(resourceEvent.UserId, resourceEvent.ResourceEventInfo["strategy"], resourceEvent.ResourceUpdateInfo)
		} else {
			log.Infof("Nothing doned for event %v", resourceEvent)
		}

	default:
		log.Infof("Nothing doned for event %v", resourceEvent)
	}

	if err != nil {
		return err
	}
	return nil
}
