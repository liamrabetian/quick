package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null"`
	Token    string `gorm:"type:varchar(255);unique;not null" json:"token"`
}
