package models

import (
	"net/http"

	//log "github.com/sirupsen/logrus"

	"github.com/108356037/algotrade/v2/auth-service/internal/database/redis"
	"github.com/108356037/algotrade/v2/auth-service/internal/jwtoken"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var bufferUser = []*User{}

type User struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
	//Role     string `json:"role"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) SignUp(c *gin.Context) {
	user := &User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	for _, existUser := range bufferUser {
		if user.Username == existUser.Username {
			c.JSON(http.StatusBadRequest, "username exists")
			return
		}
	}
	user.UserId = uuid.NewV4().String()[:8]
	bufferUser = append(bufferUser, user)

	tokenDetails, err := jwtoken.TokenPair(user.UserId, user.Username)
	if err != nil {
		c.JSON(http.StatusUnavailableForLegalReasons, err.Error())
		return
	}

	if err := jwtoken.TokenSave(user.Username, tokenDetails); err != nil {
		c.JSON(http.StatusUnavailableForLegalReasons, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	})

}

func (u *User) SignIn(c *gin.Context) {
	user := &User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	for _, existUser := range bufferUser {
		if user.Username == existUser.Username && user.Password == existUser.Password {

			existingToken := redis.Client.Keys(redis.Ctx, user.Username+"*")

			if len(existingToken.Val()) != 0 {
				for _, k := range existingToken.Val() {
					redis.Client.Del(redis.Ctx, k)
					//log.Info(delres.Result())
				}
			}

			tokenDetails, err := jwtoken.TokenPair(existUser.UserId, existUser.Username)
			if err != nil {
				c.JSON(http.StatusUnavailableForLegalReasons, err.Error())
				return
			}

			if err := jwtoken.TokenSave(user.Username, tokenDetails); err != nil {
				c.JSON(http.StatusUnavailableForLegalReasons, err.Error())
				return
			}

			c.JSON(http.StatusOK, map[string]string{
				"access_token":  tokenDetails.AccessToken,
				"refresh_token": tokenDetails.RefreshToken,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, "User not found")

}
func (u *User) SignOut(c *gin.Context)    {}
func (u *User) UserInfo(c *gin.Context)   {}
func (u *User) UserDelete(c *gin.Context) {}
