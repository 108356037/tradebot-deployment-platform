package publish

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/108356037/algotrade/v2/auth-service/global"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	AwsSession *session.Session
)

func AwsSessInit() error {
	os.Setenv("AWS_PROFILE", global.AwsSetting.Profile)
	os.Setenv("AWS_REGION", global.AwsSetting.Region)

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	AwsSession = sess
	return nil
}

// publish event when refresh token deleted
func PubRtDelete(tokenUUID string) error {
	svc := sqs.New(AwsSession)
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"event-type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("token-deleted"),
			},
		},
		MessageBody:    aws.String(tokenUUID),
		QueueUrl:       aws.String(global.AwsSetting.QueueURL),
		MessageGroupId: aws.String("token-event"),
	})
	if err != nil {
		return err
	}
	log.Infof("msg sended aws-sqs, body MD5: %s\n", *result.MD5OfMessageBody)
	return nil
}

// publish event when refresh token expired
func PubRtExpire(tokenUUID string) error {
	svc := sqs.New(AwsSession)
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"event-type": {
				DataType:    aws.String("String"),
				StringValue: aws.String("token-expired"),
			},
		},
		MessageBody:    aws.String(tokenUUID),
		QueueUrl:       aws.String(global.AwsSetting.QueueURL),
		MessageGroupId: aws.String("token-event"),
	})
	if err != nil {
		return err
	}
	log.Infof("msg sended aws-sqs, body MD5: %s\n", *result.MD5OfMessageBody)
	return nil
}
