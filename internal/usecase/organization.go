package usecase

import (
	"context"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// OrganizationUsecase 組織ユースケース
type OrganizationUsecase struct {
	organizationService *service.OrganizationService
}

// NewOrganizationUsecase 組織ユースケースを作成
func NewOrganizationUsecase(organizationService *service.OrganizationService) *OrganizationUsecase {
	return &OrganizationUsecase{
		organizationService: organizationService,
	}
}

// CreateOrganizationRequest 組織作成リクエスト
type CreateOrganizationRequest struct {
	Mail string `json:"mail" validate:"required,email"`
	Name string `json:"org_name" validate:"required"`
}

// GetOrganizations 組織一覧取得
func (u *OrganizationUsecase) GetOrganizations(ctx context.Context) ([]model.Organization, error) {
	organizations, err := u.organizationService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

// CreateOrganization 組織を作成
func (u *OrganizationUsecase) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*model.Organization, error) {
	organization, err := u.organizationService.Create(ctx, req.Mail, req.Name)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

// GetOrganization 組織を取得
func (u *OrganizationUsecase) GetOrganization(ctx context.Context, orgID string) (*model.Organization, error) {
	organization, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

// DeleteOrganization 組織を削除
func (u *OrganizationUsecase) DeleteOrganization(ctx context.Context, orgID string) error {
	if err := u.organizationService.Delete(ctx, orgID); err != nil {
		return err
	}
	return nil
}
