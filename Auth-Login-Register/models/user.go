package models

import (
	"time"
)

type User struct {
	Id                int64  `gorm:"primaryKey" json:"id"`
	NamaLengkap       string `gorm:"type:varchar(300)" json:"nama_lengkap"`
	Username          string `gorm:"type:varchar(300)" json:"username"`
	Email             string `gorm:"type:varchar(300)" json:"email"` 
	Password          string `gorm:"type:varchar(300)" json:"password"`
	IsVerified        bool   `gorm:"default:false" json:"is_verified"`            
	VerificationToken string `gorm:"type:varchar(300)" json:"verification_token"` 

	ResetToken        string    `gorm:"type:varchar(255)" json:"reset_token"`
    ResetTokenExpiry  *time.Time `json:"reset_token_expiry"`
}
