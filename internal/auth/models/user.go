package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name            string `gorm:"size:255"`
	Email           string `gorm:"size:255;unique"`
	Password        string `gorm:"size:255"`
	ResetToken      string `gorm:"size:255"`
	TokenExpiration *time.Time
}
