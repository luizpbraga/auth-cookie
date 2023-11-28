package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luizpbraga/auth-cookie/src/routes"
)

func main() {

	app := gin.Default()

	routes.Setup(app)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	app.Run()
}
