package models

import "time"

type Resource struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	ResourceType string    `gorm:"column:resource_type;not null"`
	Status       string    `gorm:"default:'active'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
