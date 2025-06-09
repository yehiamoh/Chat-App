// pkg/models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID    string `gorm:"uniqueIndex"`
	Email       string `gorm:"uniqueIndex"`
	Name        string
	Picture     string
	LastLoginAt time.Time
}