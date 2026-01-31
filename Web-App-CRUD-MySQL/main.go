package main

import (
	"net/http"
	"lapor-pak/config"      // Panggil folder config
	"lapor-pak/controllers" // Panggil folder controllers

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Nyalakan Database
	config.ConnectDB()

	// 2. Setup Router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// 3. Routing (Menghubungkan URL ke Controller)
	
	// Halaman Depan
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API Routes (Panggil fungsi dari package controllers)
	r.GET("/laporan", controllers.Index)
	r.POST("/laporan", controllers.Store)
	r.PUT("/laporan/:id", controllers.Update)
	r.DELETE("/laporan/:id", controllers.Delete)

	// 4. Jalankan
	r.Run(":8080")
}