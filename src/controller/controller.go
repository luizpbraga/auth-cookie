package controller

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/luizpbraga/auth-cookie/src/database"
	"github.com/luizpbraga/auth-cookie/src/database/models"
	"golang.org/x/crypto/bcrypt"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func Register(c *gin.Context) {
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		log.Fatal(err)
		return
	}

	name := data["name"]
	email := data["email"]
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	database.InserUserTable(&user)

	c.IndentedJSON(http.StatusOK, user)
	// c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
