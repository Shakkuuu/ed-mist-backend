package model

import (
	"time"

	"gorm.io/gorm"
)

// Organization 組織モデル
type Organization struct {
	ID        string         `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	Mail      string         `gorm:"column:mail;type:varchar(255);uniqueIndex;not null" json:"mail"`
	Name      string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}
