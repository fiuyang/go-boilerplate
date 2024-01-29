package model

import "time"

type PasswordResets struct {
	Id           int       `gorm:"type:int;primary_key"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Otp          int       `gorm:"unique;default:null"`
	CreatedAt    time.Time `gorm:"default:null"`

	User      Users     `gorm:"foreignKey:Email;references:Email"`
}
