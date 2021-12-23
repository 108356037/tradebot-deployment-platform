package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/108356037/algotrade/v2/auth-service/global"
	"github.com/108356037/algotrade/v2/auth-service/internal/database/redis"
	"github.com/108356037/algotrade/v2/auth-service/internal/event/publish"
	"github.com/108356037/algotrade/v2/auth-service/internal/event/subscribe"
	"github.com/108356037/algotrade/v2/auth-service/internal/jwtoken"
	"github.com/108356037/algotrade/v2/auth-service/internal/router"
	"github.com/108356037/algotrade/v2/auth-service/pkg"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initSetting() {
	setting, err := pkg.NewSettings()
	Fatal(err)

	err = setting.ReadSection("Server", &global.ServerSetting)
	Fatal(err)

	err = setting.ReadSection("Redis", &global.RedisSetting)
	Fatal(err)

	err = setting.ReadSection("JwtKeys", &global.JwtKeysSetting)
	Fatal(err)

	err = setting.ReadSection("Aws", &global.AwsSetting)
	Fatal(err)
}

func initComponents() {
	err := redis.RedisInit()
	Fatal(err)

	err = jwtoken.JwtInit()
	Fatal(err)

	subscribe.SubscribeInit()
	publish.AwsSessInit()
}

func main() {
	initSetting()
	initComponents()

	gin.SetMode(global.ServerSetting.RunMode)
	r := router.NewRouter()

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + global.ServerSetting.HttpPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(global.ServerSetting.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(global.ServerSetting.ReadTimeout) * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
