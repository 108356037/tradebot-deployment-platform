package redis

import (
	"context"

	"github.com/108356037/algotrade/v2/auth-service/global"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client
var Ctx = context.Background()

func RedisInit() error {
	//Initializing redis
	Client = redis.NewClient(&redis.Options{
		Addr: global.RedisSetting.Host, //redis port
	})

	//try pinging one time for connection check
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	// setting the config settings of redis
	for k, v := range global.RedisSetting.ConfigSet {
		_, err := Client.ConfigSet(Ctx, k, v).Result()
		if err != nil {
			return err
		}
	}

	return nil
}
