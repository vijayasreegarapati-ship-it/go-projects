package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	ResourceID uint      `json:"resource_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status" gorm:"default:'active'"`

	User     User     `json:"-" gorm:"foreignKey:UserID"`
	Resource Resource `json:"-" gorm:"foreignKey:ResourceID"`
}

type CreateBookingRequest struct {
	ResourceID uint      `json:"resource_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
}
