package controller

import (
	"net/http"
	"time"

	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/luizpbraga/auth-cookie/src/database"
	"github.com/luizpbraga/auth-cookie/src/database/models"
	"golang.org/x/crypto/bcrypt"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// search by email
func Login(c *gin.Context) {
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		log.Fatal(err)
		return
	}

	user, err := database.FindUserByEmail(data["email"])

	if err != nil {
		log.Fatal(err)
	}

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "rong password or email"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour + 24).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte("secret"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not login"})
		return
	}

	c.SetCookie("jwt", token, int(time.Now().Add(time.Hour+24).Unix()), "", "localhost", false, true)

	if _, err := c.Cookie("jwt"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not set the cookie"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
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
