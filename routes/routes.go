package routes

import (
	"bluebell/controllers"
	"bluebell/logger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/version", controllers.Ping)
	return r
}
