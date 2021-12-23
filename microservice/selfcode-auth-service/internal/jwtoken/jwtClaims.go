package jwtoken

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessTokenClaims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
	//TokenUUID string
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	//TokenUUID string
}
