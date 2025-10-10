package model

import (
	"time"
)

// Lesson 授業モデル（時間割の1コマ）
type Lesson struct {
	ID        string `gorm:"primaryKey;type:uuid;column:id;not null" json:"id"`
	SubjectID string `gorm:"type:uuid;column:subject_id;not null;index" json:"subject_id"`
	RoomID    string `gorm:"type:uuid;column:room_id;not null;index" json:"room_id"`
	OrgID     string `gorm:"type:uuid;column:org_id;not null;index" json:"org_id"`

	// 時間情報
	DayOfWeek int       `gorm:"column:day_of_week;not null;index" json:"day_of_week"` // 0=日, 1=月, ..., 6=土
	StartTime time.Time `gorm:"column:start_time;not null;index" json:"start_time"`   // 開始時刻
	EndTime   time.Time `gorm:"column:end_time;not null" json:"end_time"`             // 終了時刻

	// 特定日付の授業の場合（オプション）
	Date *time.Time `gorm:"column:date;index" json:"date,omitempty"`

	// 時限（オプション）
	Period int `gorm:"column:period" json:"period,omitempty"` // 1限、2限など

	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// リレーション
	Subject      Subject      `gorm:"foreignKey:SubjectID;references:ID" json:"subject,omitempty"`
	Room         Room         `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	Organization Organization `gorm:"foreignKey:OrgID;references:ID" json:"organization,omitempty"`
}

// TableName テーブル名を指定
func (Lesson) TableName() string {
	return "lessons"
}
