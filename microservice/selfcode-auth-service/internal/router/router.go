package router

import (
	"github.com/108356037/algotrade/v2/auth-service/internal/models"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	userHandlers := models.NewUser()

	apiv1 := r.Group("/api/v1/auth")
	{
		apiv1.POST("/user/signup", userHandlers.SignUp)    // for user signup
		apiv1.POST("/user/signin", userHandlers.SignIn)    // for user signin
		apiv1.POST("/user/signout", userHandlers.SignOut)  // for user signout
		apiv1.GET("/user/:id", userHandlers.UserInfo)      // get current user info
		apiv1.DELETE("/user/:id", userHandlers.UserDelete) // delete user from db
	}

	return r
}
