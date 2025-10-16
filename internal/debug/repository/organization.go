package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// OrganizationRepository 組織リポジトリ
type OrganizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository 組織リポジトリを作成
func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// FindAll 全組織を取得
func (r *OrganizationRepository) FindAll(ctx context.Context) ([]model.Organization, error) {
	var organizations []model.Organization
	err := r.db.WithContext(ctx).Find(&organizations).Error
	return organizations, err
}

// Create 組織を作成
func (r *OrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
	return r.db.WithContext(ctx).Create(organization).Error
}

// DeleteAll 全組織を削除
func (r *OrganizationRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM organizations").Error
}
