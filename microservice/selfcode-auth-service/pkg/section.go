package pkg

// import (
// 	log "github.com/sirupsen/logrus"
// )

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  int64
	WriteTimeout int64
}

type RedisSettingS struct {
	Host      string
	ConfigSet map[string]string
}

type JwtKeysSettingS struct {
	PrivKeyPath string
	PubKeyPath  string
	Issuer      string
	Subject     string
}

type AwsSettingS struct {
	Profile  string
	Region   string
	QueueID  string
	QueueURL string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
