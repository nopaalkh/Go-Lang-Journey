package controllers

import (
	"belajar-auth/config"
	"belajar-auth/models"
	"belajar-auth/utils"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Struct untuk Validasi Input
type RegisterInput struct {
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// --- HELPER: GENERATE RANDOM TOKEN (HEX) ---
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ==========================================
// 1. REGISTER SYSTEM
// ==========================================

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Validasi Password Kuat (Yang tadi kita buat)
	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- LOGIKA PENGECEKAN BARU (SPESIFIK) ---
	var countUsername int64
	var countEmail int64

	// Cek apakah Username sudah ada?
	config.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&countUsername)

	// Cek apakah Email sudah ada?
	config.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&countEmail)

	// Logika Pesan Error
	if countUsername > 0 && countEmail > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username dan Email sudah digunakan!"})
		return
	}

	if countUsername > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah dipakai, silakan pilih yang lain."})
		return
	}

	if countEmail > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar, silakan login saja."})
		return
	}
	// ------------------------------------------

	// 2. Hash Password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 3. Generate Kode Unik
	verifCode := uuid.New().String()

	// 4. Simpan User
	user := models.User{
		NamaLengkap:       input.NamaLengkap,
		Username:          input.Username,
		Email:             input.Email,
		Password:          string(hashPassword),
		IsVerified:        false,
		VerificationToken: verifCode,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal daftar user"})
		return
	}

	// 5. Kirim Email
	go func() {
		utils.SendVerificationEmail(input.Email, verifCode)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Daftar berhasil! Cek email Anda untuk verifikasi."})
}

func VerifyEmail(c *gin.Context) {
	code := c.Query("code")
	var user models.User

	// 1. Cek Token di Database
	if err := config.DB.Where("verification_token = ?", code).First(&user).Error; err != nil {
		// JIKA GAGAL: Tampilkan halaman verify.html dengan status Success=false
		c.HTML(http.StatusBadRequest, "verify.html", gin.H{
			"Success": false,
			"Message": "Invalid or expired verification link. Please request a new one.",
		})
		return
	}

	// 2. Jika Berhasil, Update User
	user.IsVerified = true
	user.VerificationToken = ""
	config.DB.Save(&user)

	// JIKA SUKSES: Tampilkan halaman verify.html dengan status Success=true
	c.HTML(http.StatusOK, "verify.html", gin.H{
		"Success": true,
	})
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password minimal 8 karakter")
	}
	// Cek Huruf Besar
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password harus mengandung minimal 1 huruf besar")
	}
	// Cek Angka
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password harus mengandung minimal 1 angka")
	}
	// Cek Simbol
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return fmt.Errorf("password harus mengandung minimal 1 simbol (!@#$)")
	}
	return nil
}

// ==========================================
// 2. LOGIN SYSTEM
// ==========================================

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/Email atau password salah"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/Email atau password salah"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun belum aktif! Cek email Anda."})
		return
	}

	// JWT Generation
	expTime := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": expTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("rahasia-kita"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login Berhasil!"})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// ==========================================
// 3. FORGOT PASSWORD SYSTEM (BARU!)
// ==========================================

// A. Proses Kirim Email Reset (POST)
func ForgotPasswordProcess(c *gin.Context) {
	var input ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Sent"})
		return
	}

	token, _ := generateToken()
	expiry := time.Now().Add(15 * time.Minute)

	// Update DB
	user.ResetToken = token
	user.ResetTokenExpiry = &expiry
	config.DB.Save(&user)

	resetLink := "http://localhost:8080/reset-password?token=" + token

	fmt.Println("Mengirim email ke:", input.Email)
	go func() {
		utils.SendResetEmail(input.Email, resetLink)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Sent"})
}

// B. Proses Update Password Baru (POST)
func ResetPasswordProcess(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Query cek token valid & belum expired
	if err := config.DB.Where("reset_token = ? AND reset_token_expiry > ?", input.Token, time.Now()).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token invalid or expired"})
		return
	}

	if err := validatePassword(input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	// Update Data & Hapus Token
	user.Password = string(hashPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// ==========================================
// 4. HTML RENDERERS
// ==========================================

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Tampilkan Halaman Forgot Password
func ShowForgotPasswordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot.html", nil)
}

// Tampilkan Halaman Reset Password (Cek token dulu)
func ShowResetPasswordPage(c *gin.Context) {
	token := c.Query("token")

	// Cek Token valid gak di DB sebelum nampilin halaman
	var user models.User
	if err := config.DB.Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Link Reset Password tidak valid atau sudah kadaluarsa.")
		return
	}

	c.HTML(http.StatusOK, "reset_password.html", nil)
}
