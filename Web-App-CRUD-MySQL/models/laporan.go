package models

// Perhatikan: package namanya 'models' (sesuai nama folder)

type Laporan struct {
	// Huruf depan field harus Besar (Exported) agar bisa dibaca package lain
	ID      uint   `json:"id" gorm:"primaryKey"`
	Pelapor string `json:"pelapor"`
	Judul   string `json:"judul"`
	Status  string `json:"status"`
}