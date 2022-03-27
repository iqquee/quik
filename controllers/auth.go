package controllers

import (
	"fmt"
	"net/http"
	"quik/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateUser(c *gin.Context) {
	var user models.User
	var wallet models.UserWallet

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Info(c.Request)

	_, userErr := models.GetUserByUsername(&user, user.Username)
	if userErr != nil {
		//hash the users password before creating user
		pwdHash, pwdErr := HashPassword(user.Password)
		if pwdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": pwdErr.Error(),
			})
			return
		}
		user.Password = string(pwdHash)

		err := models.CreateUser(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		//assign the username to the username field on the user wallet so it serves as a reference
		wallet.Username = user.Username
		//for everyuser who signes up, have their wallet set to 0  by default
		wallet.Wallet = 0
		//if the user was successfully created then create a wallet for the use
		walletErr := models.CreateWallet(&wallet)
		if walletErr != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "there was an error creating your wallet",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user was successfully created",
			"user":    user,
			"wallet":  wallet,
		})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "username already taken by another user",
		})
		return
	}
}

func SignIn(c *gin.Context) {
	var user models.Login
	var foundUser models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Info(c.Request)
	res, err := models.GetUserByUsername(&foundUser, user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the user with username %s does not exist", user.Username),
		})
		return
	}

	pwdVerifyErr := VerifyPassword(res.Password, user.Password)
	if pwdVerifyErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password mismatch, please try again",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}
