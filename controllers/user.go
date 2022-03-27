package controllers

import (
	"net/http"
	"quik/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetUsers(c *gin.Context) {
	var user []models.User
	log.Info(c.Request)
	err := models.GetAllUsers(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"users": user,
		})
	}
}
