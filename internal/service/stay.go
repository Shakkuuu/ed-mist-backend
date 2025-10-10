package service

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"
)

// StayService 滞在サービス
type StayService struct {
	stayRepo *repository.StayRepository
}

// NewStayService 滞在サービスを作成
func NewStayService(stayRepo *repository.StayRepository) *StayService {
	return &StayService{
		stayRepo: stayRepo,
	}
}

// Create 滞在を作成
func (s *StayService) Create(ctx context.Context, userID, roomID, subjectID, description string) (*model.Stay, error) {
	stay := &model.Stay{
		UserID:      userID,
		RoomID:      roomID,
		SubjectID:   subjectID,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if err := s.stayRepo.Create(ctx, stay); err != nil {
		return nil, err
	}
	return stay, nil
}

// GetAll 全滞在を取得
func (s *StayService) GetAll(ctx context.Context) ([]model.Stay, error) {
	stays, err := s.stayRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return stays, nil
}

// GetByID IDで滞在を取得
func (s *StayService) GetByID(ctx context.Context, id int) (*model.Stay, error) {
	stay, err := s.stayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return stay, nil
}

// GetByUserID ユーザーIDで滞在一覧を取得
func (s *StayService) GetByUserID(ctx context.Context, userID string) ([]model.Stay, error) {
	stays, err := s.stayRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return stays, nil
}

// GetActiveByUserID ユーザーIDでアクティブな滞在を取得
func (s *StayService) GetActiveByUserID(ctx context.Context, userID string) (*model.Stay, error) {
	stay, err := s.stayRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return stay, nil
}

// GetByRoomID 部屋IDで滞在一覧を取得
func (s *StayService) GetByRoomID(ctx context.Context, roomID string) ([]model.Stay, error) {
	stays, err := s.stayRepo.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return stays, nil
}

// GetActiveByRoomID 部屋IDでアクティブな滞在一覧を取得
func (s *StayService) GetActiveByRoomID(ctx context.Context, roomID string) ([]model.Stay, error) {
	stays, err := s.stayRepo.FindActiveByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return stays, nil
}

// GetByLessonID LessonIDで滞在一覧を取得
func (s *StayService) GetByLessonID(ctx context.Context, lessonID string) ([]model.Stay, error) {
	stays, err := s.stayRepo.FindByLessonID(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	return stays, nil
}

// GetActiveByUserAndLesson ユーザーIDとLessonIDでアクティブな滞在を取得
func (s *StayService) GetActiveByUserAndLesson(ctx context.Context, userID string, lessonID string) (*model.Stay, error) {
	stay, err := s.stayRepo.FindActiveByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}
	return stay, nil
}

// CreateWithLesson 滞在を作成（授業付き）
func (s *StayService) CreateWithLesson(ctx context.Context, stay *model.Stay) error {
	return s.stayRepo.Create(ctx, stay)
}

// Update 滞在を更新
func (s *StayService) Update(ctx context.Context, stay *model.Stay, subjectID, description string) error {
	stay.SubjectID = subjectID
	stay.Description = description

	if err := s.stayRepo.Update(ctx, stay); err != nil {
		return err
	}
	return nil
}

// EndStay 滞在を終了する
func (s *StayService) EndStay(ctx context.Context, id int) error {
	if err := s.stayRepo.EndStay(ctx, id); err != nil {
		return err
	}
	return nil
}

// Delete 滞在を削除
func (s *StayService) Delete(ctx context.Context, id int) error {
	if err := s.stayRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
