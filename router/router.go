package router

import (
	"HiChat/middleware"
	"HiChat/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")

	user := v1.Group("user")
	{
		user.GET("/list", middleware.JWY(), service.List)
		user.POST("/login_pw", middleware.JWY(), service.LoginByNameAndPassWord)
		user.POST("/new", middleware.JWY(), service.NewUser)
		user.POST("/delete", middleware.JWY(), service.DeleteUser)
		user.POST("/update", middleware.JWY(), service.UpdateUser)
	}

	return router
}
