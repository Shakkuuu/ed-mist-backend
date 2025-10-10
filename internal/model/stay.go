package model

import (
	"time"
)

// Stay 滞在モデル
type Stay struct {
	ID          int        `gorm:"primaryKey;autoIncrement;column:id;not null" json:"id"`
	UserID      string     `gorm:"type:uuid;column:user_id;not null;index" json:"user_id"`
	IsActive    bool       `gorm:"column:is_active;not null;index" json:"is_active"`
	RoomID      string     `gorm:"type:uuid;column:room_id;not null;index" json:"room_id"`
	SubjectID   string     `gorm:"type:uuid;column:subject_id;index" json:"subject_id,omitempty"`
	LessonID    *string    `gorm:"type:uuid;column:lesson_id;index" json:"lesson_id,omitempty"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty"`
	Source      string     `gorm:"column:source;type:varchar(20);default:'auto'" json:"source"` // "auto" or "manual"
	CreatedAt   time.Time  `gorm:"column:created_at;not null;index" json:"created_at"`
	LeavedAt    *time.Time `gorm:"column:leaved_at" json:"leaved_at,omitempty"`

	// リレーション
	User    User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Room    Room    `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	Subject Subject `gorm:"foreignKey:SubjectID;references:ID" json:"subject,omitempty"`
	Lesson  *Lesson `gorm:"foreignKey:LessonID;references:ID" json:"lesson,omitempty"`
}
