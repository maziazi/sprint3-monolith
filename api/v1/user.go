package v1

import (
	"github.com/gin-gonic/gin"
	"sprint3/internal/handler"
	"sprint3/internal/middleware"
)

func RegisterUserRouter(router *gin.RouterGroup) {

	{
		router.POST("/register/email", handler.RegisterUserEmail)
		router.POST("/register/phone", handler.RegisterUserPhone)
		router.POST("/login/email", handler.LoginUserEmail)
		router.POST("/login/phone", handler.LoginUserPhone)
	}
	//Kegunaan Protected itu ntar buat kalau mau akses itu harus login
	protected := router.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		/*
				USE EXAMPLE
			protected.GET("/", handler.GetUserProfileHandler)
		*/

	}
}
