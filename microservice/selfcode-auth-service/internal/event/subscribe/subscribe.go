package subscribe

import (
	"github.com/108356037/algotrade/v2/auth-service/internal/database/redis"
	"github.com/108356037/algotrade/v2/auth-service/internal/event/publish"
	log "github.com/sirupsen/logrus"
)

func SubscribeInit() {
	go SubscribeRedisKeyEvent()
}

// subscribes to data that is expired/delete/created in redis
func SubscribeRedisKeyEvent() {
	topic := redis.Client.PSubscribe(redis.Ctx, "__keyevent@0__:*")
	channel := topic.Channel()
	for {
		msg := <-channel
		switch {
		case msg.Channel == "__keyevent@0__:expire":
			log.Infof("token: %s created\n", msg.Payload)

		case msg.Channel == "__keyevent@0__:expired":
			log.Infof("token: %s expired\n", msg.Payload)
			if err := publish.PubRtExpire(msg.Payload); err != nil {
				log.Warn(err.Error())
			}

		case msg.Channel == "__keyevent@0__:del":
			log.Infof("token: %s deleted\n", msg.Payload)
			if err := publish.PubRtDelete(msg.Payload); err != nil {
				log.Warn(err.Error())
			}

		default:
			log.Warnf("no handler for this case, msg.Channel: %s\n", msg.Channel)
		}
	}
}
