package repository

import (
	"context"
	"errors"

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

// Create 部屋を作成
func (r *RoomRepository) Create(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

// FindByID IDで部屋を取得
func (r *RoomRepository) FindByID(ctx context.Context, id string) (*model.Room, error) {
	var room model.Room
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("id = ?", id).First(&room).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &room, nil
}

// FindByOrgID 組織IDで部屋一覧を取得
func (r *RoomRepository) FindByOrgID(ctx context.Context, orgID string) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("org_id = ?", orgID).Find(&rooms).Error
	return rooms, err
}

// FindByOrgRoomID 組織IDと部屋IDで部屋を取得
func (r *RoomRepository) FindByOrgRoomID(ctx context.Context, orgID, orgRoomID string) (*model.Room, error) {
	var room model.Room
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("org_id = ? AND org_room_id = ?", orgID, orgRoomID).First(&room).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &room, nil
}

// FindAll 全部屋を取得
func (r *RoomRepository) FindAll(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Find(&rooms).Error
	return rooms, err
}

// Update 部屋を更新
func (r *RoomRepository) Update(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

// Delete 部屋を削除
func (r *RoomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Room{}, "id = ?", id).Error
}
