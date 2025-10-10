package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

var ErrorRecordNotFound error = errors.New("record not found")

// UserRepository ユーザーリポジトリ
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリを作成
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create ユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID IDでユーザーを取得
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByMail メールアドレスでユーザーを取得（最初の1件のみ）
func (r *UserRepository) FindByMail(ctx context.Context, mail string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("mail = ?", mail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindAllByMail メールアドレスでユーザー全件を取得（複数組織対応）
func (r *UserRepository) FindAllByMail(ctx context.Context, mail string) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("mail = ?", mail).Find(&users).Error
	return users, err
}

// FindByOrgIDAndMail 組織IDとメールアドレスでユーザーを取得
func (r *UserRepository) FindByOrgIDAndMail(ctx context.Context, orgID, mail string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("org_id = ? AND mail = ?", orgID, mail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByOrgID 組織IDでユーザー一覧を取得
func (r *UserRepository) FindByOrgID(ctx context.Context, orgID string) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Where("org_id = ?", orgID).Find(&users).Error
	return users, err
}

// FindAll 全ユーザーを取得
func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail", "name", "created_at", "updated_at")
		}).
		Find(&users).Error
	return users, err
}

// Update ユーザーを更新
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// SoftDelete ユーザーを削除（ソフトデリート）
func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// HardDelete ユーザーを物理削除
func (r *UserRepository) HardDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&model.User{}, "id = ?", id).Error
}
