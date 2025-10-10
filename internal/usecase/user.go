package usecase

import (
	"context"
	"errors"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// UserUsecase ユーザーユースケース
type UserUsecase struct {
	userService         *service.UserService
	organizationService *service.OrganizationService
}

// NewUserUsecase ユーザーユースケースを作成
func NewUserUsecase(userService *service.UserService, organizationService *service.OrganizationService) *UserUsecase {
	return &UserUsecase{
		userService:         userService,
		organizationService: organizationService,
	}
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	OrgID    string `json:"org_id" validate:"required"`
	UserMail string `json:"user_mail" validate:"required,email"`
}

// CreateUser ユーザーを作成
func (u *UserUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	user, err := u.userService.Create(ctx, req.OrgID, req.UserMail)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsersByOrgID 組織IDでユーザー一覧を取得
func (u *UserUsecase) GetUsersByOrgID(ctx context.Context, orgID string) ([]model.User, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	users, err := u.userService.GetByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser ユーザーを取得
func (u *UserUsecase) GetUser(ctx context.Context, orgID, userID string) (*model.User, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	user, err := u.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// ユーザーが指定された組織に属しているかチェック
	if user.OrgID != orgID {
		return nil, errors.New("指定されたユーザーは組織に属していません")
	}

	return user, nil
}

// DeleteUser ユーザーを削除
func (u *UserUsecase) DeleteUser(ctx context.Context, orgID, userID string) error {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// ユーザーの存在確認
	user, err := u.userService.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// ユーザーが指定された組織に属しているかチェック
	if user.OrgID != orgID {
		return errors.New("指定されたユーザーは組織に属していません")
	}

	if err := u.userService.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}
