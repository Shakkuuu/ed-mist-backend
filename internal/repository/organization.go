package repository

import (
	"context"
	"errors"

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

// Create 組織を作成
func (r *OrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
	return r.db.WithContext(ctx).Create(organization).Error
}

// FindByID IDで組織を取得
func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*model.Organization, error) {
	var organization model.Organization
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&organization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &organization, nil
}

// FindByMail メールアドレスで組織を取得
func (r *OrganizationRepository) FindByMail(ctx context.Context, mail string) (*model.Organization, error) {
	var organization model.Organization
	err := r.db.WithContext(ctx).Where("mail = ?", mail).First(&organization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &organization, nil
}

// FindAll 全組織を取得
func (r *OrganizationRepository) FindAll(ctx context.Context) ([]model.Organization, error) {
	var organizations []model.Organization
	err := r.db.WithContext(ctx).Find(&organizations).Error
	return organizations, err
}

// Update 組織を更新
func (r *OrganizationRepository) Update(ctx context.Context, organization *model.Organization) error {
	return r.db.WithContext(ctx).Save(organization).Error
}

// SoftDelete 組織を削除（ソフトデリート）
func (r *OrganizationRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, "id = ?", id).Error
}

// HardDelete 組織を物理削除
func (r *OrganizationRepository) HardDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&model.Organization{}, "id = ?", id).Error
}
