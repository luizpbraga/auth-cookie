package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luizpbraga/auth-cookie/src/controller"
)

func Setup(r *gin.Engine) {
	r.GET("/ping", controller.Pong)
}
