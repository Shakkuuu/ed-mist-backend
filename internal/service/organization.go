package service

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// OrganizationService 組織サービス
type OrganizationService struct {
	organizationRepo *repository.OrganizationRepository
}

// NewOrganizationService 組織サービスを作成
func NewOrganizationService(organizationRepo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		organizationRepo: organizationRepo,
	}
}

// Create 組織を作成
func (o *OrganizationService) Create(ctx context.Context, mail, name string) (*model.Organization, error) {
	organization := &model.Organization{
		ID:        uuid.NewString(),
		Mail:      mail,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := o.organizationRepo.Create(ctx, organization); err != nil {
		return nil, err
	}
	return organization, nil
}

// GetAll 全組織を取得
func (o *OrganizationService) GetAll(ctx context.Context) ([]model.Organization, error) {
	organizations, err := o.organizationRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

// GetByID IDで組織を取得
func (o *OrganizationService) GetByID(ctx context.Context, id string) (*model.Organization, error) {
	organization, err := o.organizationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

// GetByMail メールアドレスで組織を取得
func (o *OrganizationService) GetByMail(ctx context.Context, mail string) (*model.Organization, error) {
	organization, err := o.organizationRepo.FindByMail(ctx, mail)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

// Update 組織を更新
func (o *OrganizationService) Update(ctx context.Context, organization *model.Organization, mail, name string) error {
	organization.Mail = mail
	organization.Name = name
	organization.UpdatedAt = time.Now()

	if err := o.organizationRepo.Update(ctx, organization); err != nil {
		return err
	}
	return nil
}

// Delete 組織を削除
func (o *OrganizationService) Delete(ctx context.Context, id string) error {
	if err := o.organizationRepo.SoftDelete(ctx, id); err != nil {
		return err
	}
	return nil
}
