package main

/// GO MODULES SUCKS!!!!!
import (
	"github.com/gin-gonic/gin"
	"github.com/luizpbraga/auth-cookie/src/database"
	"github.com/luizpbraga/auth-cookie/src/routes"
)

func main() {

	db, err := database.Connect()

	if err != nil {
		return
	}

	defer db.Close()

	app := gin.Default()

	// wtf is a midware

	routes.Setup(app)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	app.Run()
}
