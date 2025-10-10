package service

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// LessonService 授業サービス
type LessonService struct {
	lessonRepo *repository.LessonRepository
}

// NewLessonService 授業サービスを作成
func NewLessonService(lessonRepo *repository.LessonRepository) *LessonService {
	return &LessonService{
		lessonRepo: lessonRepo,
	}
}

// Create 授業を作成
func (s *LessonService) Create(ctx context.Context, subjectID, roomID, orgID string, dayOfWeek int, startTime, endTime time.Time, period int, date *time.Time) (*model.Lesson, error) {
	lesson := &model.Lesson{
		ID:        uuid.NewString(),
		SubjectID: subjectID,
		RoomID:    roomID,
		OrgID:     orgID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
		Period:    period,
		Date:      date,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.lessonRepo.Create(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

// GetByID IDで授業を取得
func (s *LessonService) GetByID(ctx context.Context, id string) (*model.Lesson, error) {
	return s.lessonRepo.FindByID(ctx, id)
}

// GetByOrgID 組織IDで授業一覧を取得
func (s *LessonService) GetByOrgID(ctx context.Context, orgID string) ([]model.Lesson, error) {
	return s.lessonRepo.FindByOrgID(ctx, orgID)
}

// GetByDate 特定の日付の授業を取得
func (s *LessonService) GetByDate(ctx context.Context, orgID string, date time.Time) ([]model.Lesson, error) {
	return s.lessonRepo.FindByDate(ctx, orgID, date)
}

// GetByUserAndDate 特定ユーザーの特定日付の授業を取得
func (s *LessonService) GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]model.Lesson, error) {
	return s.lessonRepo.FindByUserAndDate(ctx, userID, date)
}

// GetMonitoringLessons 監視対象の授業を取得
func (s *LessonService) GetMonitoringLessons(ctx context.Context, currentTime time.Time) ([]model.Lesson, error) {
	return s.lessonRepo.FindMonitoringLessons(ctx, currentTime)
}

// Update 授業を更新
func (s *LessonService) Update(ctx context.Context, lesson *model.Lesson) error {
	lesson.UpdatedAt = time.Now()
	return s.lessonRepo.Update(ctx, lesson)
}

// Delete 授業を削除
func (s *LessonService) Delete(ctx context.Context, id string) error {
	return s.lessonRepo.Delete(ctx, id)
}

// GetByRoomAndTime 部屋IDと時刻から授業を検索
func (s *LessonService) GetByRoomAndTime(ctx context.Context, roomID string, currentTime time.Time) (*model.Lesson, error) {
	return s.lessonRepo.FindByRoomAndTime(ctx, roomID, currentTime)
}
