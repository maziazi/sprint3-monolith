package v1

import (
	"github.com/gin-gonic/gin"
	"sprint3/internal/handler"
	"sprint3/internal/middleware"
)

func RegisterFileRoutes(router *gin.RouterGroup) {

	protected := router.Group("file")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/", handler.UploadFileHandler)

	}

}
