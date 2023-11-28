package routes

import (
	"../scontroller"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/ping", controller.Pong)
}
