package jwtoken

import (
	"fmt"
	"time"

	"github.com/108356037/algotrade/v2/auth-service/global"
	"github.com/108356037/algotrade/v2/auth-service/internal/database/redis"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func TokenSave(userInfo string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := redis.Client.Set(redis.Ctx, td.AccessUuid, userInfo, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := redis.Client.Set(redis.Ctx, td.RefreshUuid, userInfo, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func TokenPair(userId, username string) (*TokenDetails, error) {
	td := &TokenDetails{}

	Iat := time.Now().Unix()

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = fmt.Sprintf("%s-at-%s", username, uuid.NewV4().String()[:13])

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = fmt.Sprintf("%s-rt-%s", username, uuid.NewV4().String()[:13])

	atClaims := &AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    global.JwtKeysSetting.Issuer,
			Subject:   username,
			IssuedAt:  Iat,
			ExpiresAt: td.AtExpires,
			Id:        td.AccessUuid,
		},
		UserId:   userId,
		Username: username,
	}
	accessToken, err := AccessToken(atClaims)
	if err != nil {
		return nil, err
	}

	rtClaims := &RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    global.JwtKeysSetting.Issuer,
			IssuedAt:  Iat,
			ExpiresAt: td.RtExpires,
			Id:        td.RefreshUuid,
		},
	}
	refreshToken, err := RefreshToken(rtClaims)
	if err != nil {
		return nil, err
	}

	td.AccessToken = accessToken
	td.RefreshToken = refreshToken

	return td, nil
}

func AccessToken(claims *AccessTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func RefreshToken(claims *RefreshTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
