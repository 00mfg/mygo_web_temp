package controllers

import (
	"bluebell/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, settings.Conf.Version)
}
