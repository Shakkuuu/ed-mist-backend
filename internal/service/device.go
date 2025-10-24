package service

import (
	"context"
	"strings"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// DeviceService デバイスサービス
type DeviceService struct {
	deviceRepo *repository.DeviceRepository
}

// NewDeviceService デバイスサービスを作成
func NewDeviceService(deviceRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

// GetAll 全デバイスを取得
func (d *DeviceService) GetAll(ctx context.Context) ([]model.Device, error) {
	devices, err := d.deviceRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// GetByID IDでデバイスを取得
func (d *DeviceService) GetByID(ctx context.Context, id string) (*model.Device, error) {
	device, err := d.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// GetByDeviceID デバイスIDでデバイスを取得
func (d *DeviceService) GetByDeviceID(ctx context.Context, deviceID string) (*model.Device, error) {
	device, err := d.deviceRepo.FindByDeviceID(ctx, normalizeMACAddress(deviceID))
	if err != nil {
		return nil, err
	}
	return device, nil
}

// GetByUserID ユーザーIDでデバイス一覧を取得
func (d *DeviceService) GetByUserID(ctx context.Context, userID string) ([]model.Device, error) {
	devices, err := d.deviceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// GetActiveByUserID ユーザーIDでアクティブなデバイスを取得
func (d *DeviceService) GetActiveByUserID(ctx context.Context, userID string) (*model.Device, error) {
	device, err := d.deviceRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return device, nil
}

// normalizeMACAddress MACアドレス形式を統一（ハイフンをコロンに変換）
func normalizeMACAddress(mac string) string {
	return strings.ReplaceAll(mac, "-", ":")
}

// Create デバイスを作成
func (d *DeviceService) Create(ctx context.Context, userID, deviceID string) (*model.Device, error) {
	now := time.Now()
	device := &model.Device{
		ID:                uuid.NewString(),
		UserID:            userID,
		DeviceID:          normalizeMACAddress(deviceID), // MACアドレス形式を統一
		IsActive:          false,
		LastAuthenticated: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := d.deviceRepo.Create(ctx, device); err != nil {
		return nil, err
	}

	// Preloadを含めて再取得
	return d.deviceRepo.FindByID(ctx, device.ID)
}

// Update デバイスを更新
func (d *DeviceService) Update(ctx context.Context, device *model.Device) error {
	return d.deviceRepo.Update(ctx, device)
}

// Delete デバイスを削除
func (d *DeviceService) Delete(ctx context.Context, id string) error {
	if err := d.deviceRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// Activate デバイスをアクティブにする
func (d *DeviceService) Activate(ctx context.Context, id string) error {
	if err := d.deviceRepo.Activate(ctx, id); err != nil {
		return err
	}
	return nil
}

// ActivateWithAuthentication デバイスをアクティブにし、認証時刻を更新
func (d *DeviceService) ActivateWithAuthentication(ctx context.Context, deviceID string) (*model.Device, error) {
	device, err := d.deviceRepo.FindByDeviceID(ctx, normalizeMACAddress(deviceID))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	device.IsActive = true
	device.LastAuthenticated = now
	device.UpdatedAt = now

	if err := d.deviceRepo.Update(ctx, device); err != nil {
		return nil, err
	}

	// Preloadを含めて再取得
	return d.deviceRepo.FindByID(ctx, device.ID)
}

// Deactivate デバイスを非アクティブにする
func (d *DeviceService) Deactivate(ctx context.Context, id string) error {
	if err := d.deviceRepo.Deactivate(ctx, id); err != nil {
		return err
	}
	return nil
}

// DeactivateAllForOrg 組織の全デバイスを非アクティブにする（日次バッチ用）
func (d *DeviceService) DeactivateAllForOrg(ctx context.Context, orgID string) error {
	return d.deviceRepo.DeactivateAllForOrg(ctx, orgID)
}
