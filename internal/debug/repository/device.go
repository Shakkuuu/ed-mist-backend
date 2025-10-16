package repository

import (
	"context"

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

// FindAll 全デバイスを取得
func (r *DeviceRepository) FindAll(ctx context.Context) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

// Create デバイスを作成
func (r *DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// DeleteAll 全デバイスを削除
func (r *DeviceRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM devices").Error
}
