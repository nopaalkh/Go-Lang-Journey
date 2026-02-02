package main

import (
	"belajar-auth/config"
	"belajar-auth/controllers"
	"belajar-auth/middlewares"
	"log"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. LOAD .ENV
	err := godotenv.Load()
	if err != nil {

		log.Println("Note: File .env tidak ditemukan, menggunakan Environment System")
	}

	config.ConnectDatabase()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Rute Publik
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.ShowRegisterPage)
	r.POST("/register", controllers.Register)
	r.GET("/verify", controllers.VerifyEmail)
	r.GET("/logout", controllers.Logout)

	r.GET("/forgot-password", controllers.ShowForgotPasswordPage)
	r.POST("/forgot-password", controllers.ForgotPasswordProcess)

	// Reset Password Routes
	r.GET("/reset-password", controllers.ShowResetPasswordPage)
	r.POST("/reset-password", controllers.ResetPasswordProcess)

	// Rute Private
	private := r.Group("/")
	private.Use(middlewares.AuthMiddleware())
	{
		private.GET("/", controllers.Index)
		private.GET("/dashboard", controllers.Dashboard)
	}

	r.NoRoute(func(c *gin.Context) {
		// Tampilkan halaman 404.html dengan status 404
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// 2. PORT DINAMIS
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
