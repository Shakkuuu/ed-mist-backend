package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// RoomRepository 部屋リポジトリ
type RoomRepository struct {
	db *gorm.DB
}

// NewRoomRepository 部屋リポジトリを作成
func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// FindAll 全部屋を取得
func (r *RoomRepository) FindAll(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	return rooms, err
}

// Create 部屋を作成
func (r *RoomRepository) Create(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

// DeleteAll 全部屋を削除
func (r *RoomRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM rooms").Error
}
