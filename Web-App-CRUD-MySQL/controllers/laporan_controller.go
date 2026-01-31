package controllers

import (
	"net/http"
	"lapor-pak/config" // Import database
	"lapor-pak/models" // Import struktur tabel

	"github.com/gin-gonic/gin"
)

// 1. GET (Tampilkan Semua)
func Index(c *gin.Context) {
	var laporan []models.Laporan
	// Pakai config.DB untuk akses database
	config.DB.Order("id desc").Find(&laporan)
	c.JSON(http.StatusOK, gin.H{"data": laporan})
}

// 2. POST (Tambah Data)
func Store(c *gin.Context) {
	var input models.Laporan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status == "" {
		input.Status = "Pending"
	}

	config.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"pesan": "Berhasil", "data": input})
}

// 3. PUT (Update Status)
func Update(c *gin.Context) {
	id := c.Param("id")
	var laporan models.Laporan

	if err := config.DB.First(&laporan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	config.DB.Model(&laporan).Update("status", "Selesai")
	c.JSON(http.StatusOK, gin.H{"pesan": "Status berhasil diupdate!"})
}

// 4. DELETE (Hapus)
func Delete(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Unscoped().Delete(&models.Laporan{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pesan": "Data berhasil dihapus"})
}