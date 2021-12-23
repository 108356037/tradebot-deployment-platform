package global

import (
	"github.com/108356037/algotrade/v2/auth-service/pkg"
)

var (
	ServerSetting  *pkg.ServerSettingS
	RedisSetting   *pkg.RedisSettingS
	JwtKeysSetting *pkg.JwtKeysSettingS
	AwsSetting     *pkg.AwsSettingS
)
