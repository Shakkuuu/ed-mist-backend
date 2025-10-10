package model

import (
	"time"
)

// Zone ゾーンモデル
type Zone struct {
	ID         string    `gorm:"primaryKey;type:uuid;column:id" json:"id"`
	MistZoneID string    `gorm:"column:mist_zone_id;type:varchar(255);uniqueIndex" json:"mist_zone_id"`
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name"`
	MapID      string    `gorm:"column:map_id;type:varchar(255)" json:"map_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsActive   bool      `gorm:"column:is_active" json:"is_active"`

	// リレーション
	Map Map `gorm:"foreignKey:MapID" json:"map,omitempty"`
}
