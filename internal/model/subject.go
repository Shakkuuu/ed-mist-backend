package model

import (
	"time"
)

// Subject 科目モデル
type Subject struct {
	ID        string    `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Year      int       `gorm:"column:year;not null;index" json:"year"`
	OrgID     string    `gorm:"type:uuid;column:org_id;not null;index" json:"org_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// リレーション
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
}
