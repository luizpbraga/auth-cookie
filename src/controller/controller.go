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

func Logout(c *gin.Context) {
	if _, err := c.Cookie("jwt"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not connected"})
		return
	}

	c.SetCookie("jwt", "", -1, "", "0.0.0.0", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}

// search by email
func Login(c *gin.Context) {
	//
	// if _, err := c.Cookie("jwt"); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"message": "user already connected"})
	// 	return
	// }

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
		ExpiresAt: time.Now().Add(time.Second * 60).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte("secret"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not login"})
		return
	}

	c.SetCookie("jwt", token, int(time.Now().Add(time.Second*60).Unix()), "/", "0.0.0.0", false, true)

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
}

// vai olhar pro cookie atual e lancar um user
func GetUser(c *gin.Context) {
	cookie, err := c.Cookie("jwt")

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User Not connected or cookie not found"})
		return
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})

	// if err != nil {
	// 	c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "StatusUnauthorized"})
	// 	return
	// }

	claims := token.Claims.(*jwt.StandardClaims)

	user, err := database.FindUserById(claims.Issuer)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "StatusUnauthorized"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
