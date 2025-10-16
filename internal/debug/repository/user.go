package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// UserRepository ユーザーリポジトリ
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリを作成
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll 全ユーザーを取得
func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// Create ユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// DeleteAll 全ユーザーを削除
func (r *UserRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM users").Error
}
