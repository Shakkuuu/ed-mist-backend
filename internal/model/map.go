package model

import (
	"time"
)

// Map マップモデル
type Map struct {
	ID           string    `gorm:"primaryKey;type:uuid;column:id" json:"id"`
	MistMapID    string    `gorm:"column:mist_map_id;type:varchar(255);uniqueIndex" json:"mist_map_id"`
	Name         string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Width        float64   `gorm:"column:width" json:"width"`
	Height       float64   `gorm:"column:height" json:"height"`
	WidthM       float64   `gorm:"column:width_m" json:"width_m"`
	HeightM      float64   `gorm:"column:height_m" json:"height_m"`
	PPM          float64   `gorm:"column:ppm" json:"ppm"`
	URL          string    `gorm:"column:url;type:text" json:"url"`
	ThumbnailURL string    `gorm:"column:thumbnail_url;type:text" json:"thumbnail_url"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsActive     bool      `gorm:"column:is_active" json:"is_active"`

	// リレーション
	Zones []Zone `gorm:"foreignKey:MapID" json:"zones,omitempty"`
}
