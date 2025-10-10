package model

import (
	"time"

	"gorm.io/gorm"
)

// User ユーザーモデル
type User struct {
	ID        string         `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	OrgID     string         `gorm:"type:uuid;column:org_id;not null;index" json:"org_id"`
	Mail      string         `gorm:"column:mail;type:varchar(255);not null;uniqueIndex" json:"mail"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	// リレーション
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
	Devices      []Device     `gorm:"foreignKey:UserID" json:"devices,omitempty"`
}
