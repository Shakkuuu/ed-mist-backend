package service

import (
	"context"
	"errors"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// UserService ユーザーサービス
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService ユーザーサービスを作成
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Create ユーザーを作成
func (u *UserService) Create(ctx context.Context, orgID, mail string) (*model.User, error) {
	user := &model.User{
		ID:        uuid.NewString(),
		OrgID:     orgID,
		Mail:      mail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Preloadを含めて再取得
	return u.userRepo.FindByID(ctx, user.ID)
}

// GetAll 全ユーザーを取得
func (u *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID IDでユーザーを取得
func (u *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// GetByMail メールアドレスでユーザーを取得（最初の1件のみ）
func (u *UserService) GetByMail(ctx context.Context, mail string) (*model.User, error) {
	user, err := u.userRepo.FindByMail(ctx, mail)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// GetAllByMail メールアドレスでユーザー全件を取得（複数組織対応）
func (u *UserService) GetAllByMail(ctx context.Context, mail string) ([]model.User, error) {
	users, err := u.userRepo.FindAllByMail(ctx, mail)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByOrgIDAndMail 組織IDとメールアドレスでユーザーを取得
func (u *UserService) GetByOrgIDAndMail(ctx context.Context, orgID, mail string) (*model.User, error) {
	user, err := u.userRepo.FindByOrgIDAndMail(ctx, orgID, mail)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// GetByOrgID 組織IDでユーザー一覧を取得
func (u *UserService) GetByOrgID(ctx context.Context, orgID string) ([]model.User, error) {
	users, err := u.userRepo.FindByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update ユーザーを更新
func (u *UserService) Update(ctx context.Context, user *model.User, mail string) error {
	user.Mail = mail
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

// Delete ユーザーを削除
func (u *UserService) Delete(ctx context.Context, id string) error {
	if err := u.userRepo.SoftDelete(ctx, id); err != nil {
		return err
	}
	return nil
}
