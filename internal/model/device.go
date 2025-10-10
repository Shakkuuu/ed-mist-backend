package model

import (
	"time"
)

// Device デバイスモデル
type Device struct {
	ID                string    `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	UserID            string    `gorm:"type:uuid;column:user_id;not null" json:"user_id"`
	DeviceID          string    `gorm:"column:device_id;type:varchar(255);not null;uniqueIndex" json:"device_id"`
	IsActive          bool      `gorm:"column:is_active;default:false" json:"is_active"`
	LastAuthenticated time.Time `gorm:"column:last_authenticated" json:"last_authenticated"`
	CreatedAt         time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// リレーション
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// IsAuthenticatedToday 本日認証済みかチェック
func (d *Device) IsAuthenticatedToday() bool {
	if !d.IsActive {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	authDay := d.LastAuthenticated.Truncate(24 * time.Hour)

	return authDay.Equal(today) || authDay.After(today)
}
