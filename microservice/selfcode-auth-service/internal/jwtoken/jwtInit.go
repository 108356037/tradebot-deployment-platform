package jwtoken

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/108356037/algotrade/v2/auth-service/global"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func JwtInit() error {
	signKeyB, err := ioutil.ReadFile(global.JwtKeysSetting.PrivKeyPath)
	if err != nil {
		return err
	}

	verifyKeyB, err := ioutil.ReadFile(global.JwtKeysSetting.PubKeyPath)
	if err != nil {
		return err
	}

	_signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signKeyB)
	if err != nil {
		return err
	}

	_verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyKeyB)
	if err != nil {
		return err
	}

	signKey = _signKey
	verifyKey = _verifyKey

	return nil
}
