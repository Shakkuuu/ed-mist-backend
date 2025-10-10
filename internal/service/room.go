package service

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// RoomService 部屋サービス
type RoomService struct {
	roomRepo *repository.RoomRepository
}

// NewRoomService 部屋サービスを作成
func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

// Create 部屋を作成
func (r *RoomService) Create(ctx context.Context, orgID, orgRoomID, name, caption, mistZoneID string) (*model.Room, error) {
	room := &model.Room{
		ID:         uuid.NewString(),
		OrgID:      orgID,
		OrgRoomID:  orgRoomID,
		Name:       name,
		Caption:    caption,
		MistZoneID: mistZoneID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := r.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	// Preloadを含めて再取得
	return r.roomRepo.FindByID(ctx, room.ID)
}

// GetAll 全部屋を取得
func (r *RoomService) GetAll(ctx context.Context) ([]model.Room, error) {
	rooms, err := r.roomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetByID IDで部屋を取得
func (r *RoomService) GetByID(ctx context.Context, id string) (*model.Room, error) {
	room, err := r.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetByOrgID 組織IDで部屋一覧を取得
func (r *RoomService) GetByOrgID(ctx context.Context, orgID string) ([]model.Room, error) {
	rooms, err := r.roomRepo.FindByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetByOrgRoomID 組織IDと部屋IDで部屋を取得
func (r *RoomService) GetByOrgRoomID(ctx context.Context, orgID, orgRoomID string) (*model.Room, error) {
	room, err := r.roomRepo.FindByOrgRoomID(ctx, orgID, orgRoomID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// Update 部屋を更新
func (r *RoomService) Update(ctx context.Context, room *model.Room, name, caption, mistZoneID string) error {
	room.Name = name
	room.Caption = caption
	room.MistZoneID = mistZoneID
	room.UpdatedAt = time.Now()

	if err := r.roomRepo.Update(ctx, room); err != nil {
		return err
	}
	return nil
}

// Delete 部屋を削除
func (r *RoomService) Delete(ctx context.Context, id string) error {
	if err := r.roomRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
