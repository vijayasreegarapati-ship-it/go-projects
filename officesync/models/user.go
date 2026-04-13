package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Role      string    `gorm:"default:'employee'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
