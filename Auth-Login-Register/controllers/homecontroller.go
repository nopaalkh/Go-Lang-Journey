package controllers

import (
	"belajar-auth/config"
	"belajar-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": user,
	})
}

func Dashboard(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": user,
	})
}
