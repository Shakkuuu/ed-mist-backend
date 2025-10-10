package usecase

import (
	"context"
	"errors"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// RoomUsecase 部屋ユースケース
type RoomUsecase struct {
	roomService         *service.RoomService
	organizationService *service.OrganizationService
}

// NewRoomUsecase 部屋ユースケースを作成
func NewRoomUsecase(roomService *service.RoomService, organizationService *service.OrganizationService) *RoomUsecase {
	return &RoomUsecase{
		roomService:         roomService,
		organizationService: organizationService,
	}
}

// CreateRoomRequest 部屋作成リクエスト
type CreateRoomRequest struct {
	OrgID      string `json:"org_id" validate:"required"`
	OrgRoomID  string `json:"org_room_id" validate:"required"`
	RoomName   string `json:"room_name" validate:"required"`
	Caption    string `json:"caption"`
	MistZoneID string `json:"mist_zone_id"`
}

// CreateRoom 部屋を作成
func (u *RoomUsecase) CreateRoom(ctx context.Context, req *CreateRoomRequest) (*model.Room, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	room, err := u.roomService.Create(ctx, req.OrgID, req.OrgRoomID, req.RoomName, req.Caption, req.MistZoneID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetRoomsByOrgID 組織IDで部屋一覧を取得
func (u *RoomUsecase) GetRoomsByOrgID(ctx context.Context, orgID string) ([]model.Room, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	rooms, err := u.roomService.GetByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// UpdateRoomRequest 部屋更新リクエスト
type UpdateRoomRequest struct {
	RoomName   string `json:"room_name"`
	Caption    string `json:"caption"`
	MistZoneID string `json:"mist_zone_id"`
}

// UpdateRoom 部屋を更新
func (u *RoomUsecase) UpdateRoom(ctx context.Context, orgID, roomID string, req *UpdateRoomRequest) (*model.Room, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// 部屋の存在確認
	room, err := u.roomService.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	// 部屋が指定された組織に属しているかチェック
	if room.OrgID != orgID {
		return nil, errors.New("指定された部屋は組織に属していません")
	}

	// 部屋を更新
	if err := u.roomService.Update(ctx, room, req.RoomName, req.Caption, req.MistZoneID); err != nil {
		return nil, err
	}

	// 更新後の部屋情報を取得
	updatedRoom, err := u.roomService.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}

// DeleteRoom 部屋を削除
func (u *RoomUsecase) DeleteRoom(ctx context.Context, orgID, roomID string) error {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// 部屋の存在確認
	room, err := u.roomService.GetByID(ctx, roomID)
	if err != nil {
		return err
	}

	// 部屋が指定された組織に属しているかチェック
	if room.OrgID != orgID {
		return errors.New("指定された部屋は組織に属していません")
	}

	if err := u.roomService.Delete(ctx, roomID); err != nil {
		return err
	}
	return nil
}
