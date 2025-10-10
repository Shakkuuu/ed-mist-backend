package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// DeviceRepository デバイスリポジトリ
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository デバイスリポジトリを作成
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create デバイスを作成
func (r *DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// FindByID IDでデバイスを取得
func (r *DeviceRepository) FindByID(ctx context.Context, id string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Where("id = ?", id).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &device, nil
}

// FindByDeviceID デバイスIDでデバイスを取得
func (r *DeviceRepository) FindByDeviceID(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &device, nil
}

// FindByUserID ユーザーIDでデバイス一覧を取得
func (r *DeviceRepository) FindByUserID(ctx context.Context, userID string) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail")
		}).
		Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

// GetActiveByUserID ユーザーIDでアクティブなデバイスを取得
func (r *DeviceRepository) GetActiveByUserID(ctx context.Context, userID string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail")
		}).
		Where("user_id = ? AND is_active = ?", userID, true).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &device, nil
}

// FindAll 全デバイスを取得
func (r *DeviceRepository) FindAll(ctx context.Context) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

// Update デバイスを更新
func (r *DeviceRepository) Update(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// Delete デバイスを削除
func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Device{}, "id = ?", id).Error
}

// Activate デバイスをアクティブにする
func (r *DeviceRepository) Activate(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", id).Update("is_active", true).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorRecordNotFound
		}
		return err
	}
	return nil
}

// Deactivate デバイスを非アクティブにする
func (r *DeviceRepository) Deactivate(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorRecordNotFound
		}
		return err
	}
	return nil
}
