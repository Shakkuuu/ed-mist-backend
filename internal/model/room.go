package model

import (
	"time"
)

// Room 部屋モデル
type Room struct {
	ID         string    `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	OrgID      string    `gorm:"type:uuid;column:org_id;not null;index" json:"org_id"`
	OrgRoomID  string    `gorm:"column:org_room_id;type:varchar(255);not null;index" json:"org_room_id"`
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name,omitempty"`
	Caption    string    `gorm:"column:caption;type:text" json:"caption,omitempty"`
	MistZoneID string    `gorm:"column:mist_zone_id;type:varchar(255);index" json:"mist_zone_id,omitempty"`
	MapID      string    `gorm:"column:map_id;type:varchar(255);index" json:"map_id,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// リレーション
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
}
