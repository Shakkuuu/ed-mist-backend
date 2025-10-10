package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"
	"github.com/Shakkuuu/ed-mist-backend/internal/service"
)

// StayLogUsecase 滞在ログユースケース
type StayLogUsecase struct {
	stayService         *service.StayService
	userService         *service.UserService
	roomService         *service.RoomService
	subjectService      *service.SubjectService
	organizationService *service.OrganizationService
}

// NewStayLogUsecase 滞在ログユースケースを作成
func NewStayLogUsecase(stayService *service.StayService, userService *service.UserService, roomService *service.RoomService, subjectService *service.SubjectService, organizationService *service.OrganizationService) *StayLogUsecase {
	return &StayLogUsecase{
		stayService:         stayService,
		userService:         userService,
		roomService:         roomService,
		subjectService:      subjectService,
		organizationService: organizationService,
	}
}

// GetStayLogsRequest 滞在ログ取得リクエスト
type GetStayLogsRequest struct {
	OrgID     string     `json:"org_id"`
	RoomID    string     `json:"room_id"`
	SubjectID string     `json:"subject_id"`
	UserID    string     `json:"user_id"`
	IsActive  *bool      `json:"is_active"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// GetStayLogs 滞在ログを取得
func (u *StayLogUsecase) GetStayLogs(ctx context.Context, req *GetStayLogsRequest) ([]model.Stay, error) {
	// 組織の存在確認
	_, err := u.organizationService.GetByID(ctx, req.OrgID)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, errors.New("組織が見つかりません")
		}
		return nil, err
	}

	// 部屋の存在確認
	room, err := u.roomService.GetByID(ctx, req.RoomID)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return nil, errors.New("部屋が見つかりません")
		}
		return nil, err
	}

	// 部屋が指定された組織に属しているかチェック
	if room.OrgID != req.OrgID {
		return nil, errors.New("指定された部屋は組織に属していません")
	}

	// 科目の存在確認（指定されている場合）
	if req.SubjectID != "" {
		subject, err := u.subjectService.GetByID(ctx, req.SubjectID)
		if err != nil {
			if errors.Is(err, repository.ErrorRecordNotFound) {
				return nil, errors.New("科目が見つかりません")
			}
			return nil, err
		}

		// 科目が指定された組織に属しているかチェック
		if subject.OrgID != req.OrgID {
			return nil, errors.New("指定された科目は組織に属していません")
		}
	}

	// ユーザーの存在確認（指定されている場合）
	if req.UserID != "" {
		user, err := u.userService.GetByID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrorRecordNotFound) {
				return nil, errors.New("ユーザーが見つかりません")
			}
			return nil, err
		}

		// ユーザーが指定された組織に属しているかチェック
		if user.OrgID != req.OrgID {
			return nil, errors.New("指定されたユーザーは組織に属していません")
		}
	}

	// 滞在ログを取得
	stays, err := u.getStayLogsWithFilters(ctx, req)
	if err != nil {
		return nil, err
	}

	return stays, nil
}

// getStayLogsWithFilters フィルター付きで滞在ログを取得
func (u *StayLogUsecase) getStayLogsWithFilters(ctx context.Context, req *GetStayLogsRequest) ([]model.Stay, error) {
	// 基本的には部屋と科目でフィルタリング
	if req.UserID != "" {
		// 特定ユーザーの滞在ログを取得
		userStays, err := u.stayService.GetByUserID(ctx, req.UserID)
		if err != nil {
			return nil, err
		}

		// フィルタリング
		var filteredStays []model.Stay
		for _, stay := range userStays {
			if u.matchesFilters(stay, req) {
				filteredStays = append(filteredStays, stay)
			}
		}
		return filteredStays, nil
	} else {
		// 部屋の滞在ログを取得
		roomStays, err := u.stayService.GetByRoomID(ctx, req.RoomID)
		if err != nil {
			return nil, err
		}

		// フィルタリング
		var filteredStays []model.Stay
		for _, stay := range roomStays {
			if u.matchesFilters(stay, req) {
				filteredStays = append(filteredStays, stay)
			}
		}
		return filteredStays, nil
	}
}

// matchesFilters フィルター条件にマッチするかチェック
func (u *StayLogUsecase) matchesFilters(stay model.Stay, req *GetStayLogsRequest) bool {
	// 科目IDフィルター
	if req.SubjectID != "" && stay.SubjectID != req.SubjectID {
		return false
	}

	// アクティブ状態フィルター
	if req.IsActive != nil && stay.IsActive != *req.IsActive {
		return false
	}

	// 時間範囲フィルター
	if req.StartTime != nil && stay.CreatedAt.Before(*req.StartTime) {
		return false
	}

	if req.EndTime != nil {
		if stay.LeavedAt == nil {
			// 退室時間が設定されていない場合は現在時刻と比較
			if time.Now().After(*req.EndTime) {
				return false
			}
		} else {
			// 退室時間が設定されている場合は退室時間と比較
			if stay.LeavedAt.After(*req.EndTime) {
				return false
			}
		}
	}

	return true
}

// ParseStayLogsQuery クエリパラメータをパース
func ParseStayLogsQuery(userID string, isActiveStr string, startTimeStr string, endTimeStr string) (*GetStayLogsRequest, error) {
	req := &GetStayLogsRequest{
		UserID: userID,
	}

	// IsActiveパラメータの解析
	if isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			return nil, errors.New("is_activeパラメータが無効です")
		}
		req.IsActive = &isActive
	}

	// 時間範囲パラメータの解析
	if startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, errors.New("start_timeパラメータが無効です")
		}
		req.StartTime = &startTime
	}

	if endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, errors.New("end_timeパラメータが無効です")
		}
		req.EndTime = &endTime
	}

	return req, nil
}
