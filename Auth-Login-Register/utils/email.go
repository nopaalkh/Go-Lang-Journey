package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

const (
	SMTP_HOST = "smtp.gmail.com"
	SMTP_PORT = 587
)

func getEmailConfig() (string, string) {
	emailAuth := os.Getenv("EMAIL_AUTH")
	emailPass := os.Getenv("EMAIL_PASSWORD")

	if emailAuth == "" || emailPass == "" {
		fmt.Println("⚠️  WARNING: Konfigurasi Email di .env kosong atau belum terbaca!")
	}

	return emailAuth, emailPass
}

// 1. KIRIM VERIFIKASI EMAIL
func SendVerificationEmail(toEmail, verificationCode string) error {
	emailAuth, emailPass := getEmailConfig()
	baseURL := os.Getenv("BASE_URL")

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("Account Security <%s>", emailAuth))
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Aktivasi Akun Baru Anda")

	link := fmt.Sprintf("%s/verify?code=%s", baseURL, verificationCode)

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; color: #333;">
			<h3>Halo! Selamat datang.</h3>
			<p>Silakan klik link di bawah ini untuk mengaktifkan akun Anda:</p>
			<a href="%s" style="background: #0d6efd; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Verifikasi Akun</a>
		</div>
	`, link)

	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(SMTP_HOST, SMTP_PORT, emailAuth, emailPass)

	return dialer.DialAndSend(mailer)
}

// 2. KIRIM RESET PASSWORD
func SendResetEmail(toEmail, link string) error {

	emailAuth, emailPass := getEmailConfig()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("App Security <%s>", emailAuth))
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Reset Password Request")

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; color: #333;">
			<h2>Reset Password</h2>
			<p>Klik tombol di bawah untuk membuat password baru:</p>
			<a href="%s" style="background: #6366f1; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
			<p style="margin-top:20px; font-size: 12px; color: #777;">Link berlaku 15 menit.</p>
		</div>
	`, link)

	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(SMTP_HOST, SMTP_PORT, emailAuth, emailPass)

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("❌ ERROR: Gagal kirim email reset:", err)
		return err
	} else {
		fmt.Println("✅ SUKSES: Email reset terkirim ke:", toEmail)
		return nil
	}
}
