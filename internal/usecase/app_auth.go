package usecase

import (
	"context"
	"errors"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

var (
	ErrorAlreadyRegistered    = errors.New("すでに登録済みです")
	ErrorUserNotFound         = errors.New("ユーザーが見つかりません")
	ErrorOrganizationNotFound = errors.New("組織が見つかりません")
)

// AppAuthUsecase アプリ認証ユースケース
type AppAuthUsecase struct {
	userService         *service.UserService
	deviceService       *service.DeviceService
	organizationService *service.OrganizationService
}

// NewAppAuthUsecase アプリ認証ユースケースを作成
func NewAppAuthUsecase(userService *service.UserService, deviceService *service.DeviceService, organizationService *service.OrganizationService) *AppAuthUsecase {
	return &AppAuthUsecase{
		userService:         userService,
		deviceService:       deviceService,
		organizationService: organizationService,
	}
}

// DeviceRegisterRequest デバイス登録リクエスト
type DeviceRegisterRequest struct {
	OrgID    *string `json:"org_id"` // オプショナル：複数組織がある場合に指定
	Mail     string  `json:"mail" validate:"required,email"`
	DeviceID string  `json:"device_id" validate:"required"`
}

// DeviceRegisterResponse デバイス登録レスポンス
type DeviceRegisterResponse struct {
	User   *model.User   `json:"user,omitempty"`
	Device *model.Device `json:"device,omitempty"`
}

// OrganizationChoice 組織選択情報
type OrganizationChoice struct {
	OrgID   string `json:"org_id"`
	OrgName string `json:"org_name"`
	OrgMail string `json:"org_mail"`
}

// MultipleOrganizationsResponse 複数組織が見つかった場合のレスポンス
type MultipleOrganizationsResponse struct {
	Message       string               `json:"message"`
	Organizations []OrganizationChoice `json:"organizations"`
	RequiresOrgID bool                 `json:"requires_org_id"`
}

// DeviceRegister デバイス登録（管理画面で作成済みのユーザーが前提）
func (u *AppAuthUsecase) DeviceRegister(ctx context.Context, req *DeviceRegisterRequest) (interface{}, error) {
	// メールアドレスで全ユーザーを取得
	users, err := u.userService.GetAllByMail(ctx, req.Mail)
	if err != nil {
		return nil, err
	}

	// ユーザーが見つからない場合
	if len(users) == 0 {
		return nil, ErrorUserNotFound
	}

	// 組織IDが指定されている場合は、その組織のユーザーを使用
	if req.OrgID != nil && *req.OrgID != "" {
		for _, user := range users {
			if user.OrgID == *req.OrgID {
				// デバイス登録処理
				device, err := u.registerDevice(ctx, user.ID, req.DeviceID)
				if err != nil {
					return nil, err
				}

				return &DeviceRegisterResponse{
					User:   &user,
					Device: device,
				}, nil
			}
		}
		// 指定された組織IDのユーザーが見つからない
		return nil, ErrorUserNotFound
	}

	// 1つだけ見つかった場合は自動的にデバイス登録
	if len(users) == 1 {
		user := users[0]
		device, err := u.registerDevice(ctx, user.ID, req.DeviceID)
		if err != nil {
			return nil, err
		}

		return &DeviceRegisterResponse{
			User:   &user,
			Device: device,
		}, nil
	}

	// 複数見つかった場合は組織選択を促す
	var organizations []OrganizationChoice
	for _, user := range users {
		organizations = append(organizations, OrganizationChoice{
			OrgID:   user.Organization.ID,
			OrgName: user.Organization.Name,
			OrgMail: user.Organization.Mail,
		})
	}

	return &MultipleOrganizationsResponse{
		Message:       "複数の組織に同じメールアドレスのユーザーが存在します。組織を選択してください。",
		Organizations: organizations,
		RequiresOrgID: true,
	}, nil
}

// registerDevice デバイス登録
func (u *AppAuthUsecase) registerDevice(ctx context.Context, userID, deviceID string) (*model.Device, error) {
	// ユーザーの既存デバイスを確認
	devices, err := u.deviceService.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 同じデバイスIDが既に存在するかチェック
	for _, device := range devices {
		if device.DeviceID == deviceID {
			if device.IsActive {
				return nil, ErrorAlreadyRegistered
			}
			// 非アクティブなデバイスを再アクティブ化
			if err := u.deviceService.Activate(ctx, device.ID); err != nil {
				return nil, err
			}
			return &device, nil
		}
	}

	// 他のユーザーが同じデバイスIDを使用している場合は非アクティブ化
	existingDevice, err := u.deviceService.GetByDeviceID(ctx, deviceID)
	if err == nil && existingDevice != nil {
		if err := u.deviceService.Deactivate(ctx, existingDevice.ID); err != nil {
			return nil, err
		}
	}

	// ユーザーの既存アクティブデバイスを非アクティブ化
	if activeDevice, err := u.deviceService.GetActiveByUserID(ctx, userID); err == nil && activeDevice != nil {
		if err := u.deviceService.Deactivate(ctx, activeDevice.ID); err != nil {
			return nil, err
		}
	}

	// 新しいデバイスを作成
	device, err := u.deviceService.Create(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}

	return device, nil
}
