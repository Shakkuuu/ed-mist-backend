package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// LessonRepository 授業リポジトリ
type LessonRepository struct {
	db *gorm.DB
}

// NewLessonRepository 授業リポジトリを作成
func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create 授業を作成
func (r *LessonRepository) Create(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

// FindByID IDで授業を取得
func (r *LessonRepository) FindByID(ctx context.Context, id string) (*model.Lesson, error) {
	var lesson model.Lesson
	err := r.db.WithContext(ctx).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year", "org_id")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "caption", "mist_zone_id")
		}).
		Where("id = ?", id).
		First(&lesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &lesson, nil
}

// FindByOrgID 組織IDで授業一覧を取得
func (r *LessonRepository) FindByOrgID(ctx context.Context, orgID string) ([]model.Lesson, error) {
	var lessons []model.Lesson
	err := r.db.WithContext(ctx).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("org_id = ?", orgID).
		Order("start_time ASC").
		Find(&lessons).Error
	return lessons, err
}

// FindByDate 特定の日付の授業を取得
func (r *LessonRepository) FindByDate(ctx context.Context, orgID string, date time.Time) ([]model.Lesson, error) {
	var lessons []model.Lesson

	dayOfWeek := int(date.Weekday())
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("org_id = ?", orgID).
		Where("day_of_week = ?", dayOfWeek).
		Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Order("start_time ASC").
		Find(&lessons).Error

	return lessons, err
}

// FindByUserAndDate 特定ユーザーの特定日付の授業を取得
// TODO: ユーザーと授業の紐付けテーブルが必要な場合はここを拡張
func (r *LessonRepository) FindByUserAndDate(ctx context.Context, userID string, date time.Time) ([]model.Lesson, error) {
	// 現状は組織の全授業を返す（ユーザーとLessonの中間テーブルがないため）
	// 実装では、ユーザーの組織IDを取得してその組織の授業を返す

	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}

	return r.FindByDate(ctx, user.OrgID, date)
}

// FindMonitoringLessons 監視対象の授業を取得（開始5分前〜終了10分後）
func (r *LessonRepository) FindMonitoringLessons(ctx context.Context, currentTime time.Time) ([]model.Lesson, error) {
	var lessons []model.Lesson

	// 現在時刻の5分後から監視終了までの授業を取得
	monitorStart := currentTime.Add(-5 * time.Minute)
	monitorEnd := currentTime.Add(10 * time.Minute)

	dayOfWeek := int(currentTime.Weekday())

	err := r.db.WithContext(ctx).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("day_of_week = ?", dayOfWeek).
		Where(
			"(start_time <= ? AND end_time >= ?) OR (start_time >= ? AND start_time <= ?)",
			monitorEnd, monitorStart, monitorStart, monitorEnd,
		).
		Find(&lessons).Error

	return lessons, err
}

// Update 授業を更新
func (r *LessonRepository) Update(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Save(lesson).Error
}

// Delete 授業を削除
func (r *LessonRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Lesson{}).Error
}

// FindByRoomAndTime 部屋IDと時刻から授業を検索
func (r *LessonRepository) FindByRoomAndTime(ctx context.Context, roomID string, currentTime time.Time) (*model.Lesson, error) {
	var lesson model.Lesson

	dayOfWeek := int(currentTime.Weekday())

	// デバッグログ
	log.Printf("[FindByRoomAndTime] 検索条件: RoomID=%s, CurrentTime=%s, DayOfWeek=%d",
		roomID, currentTime.Format("2006-01-02 15:04:05"), dayOfWeek)

	err := r.db.WithContext(ctx).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("room_id = ?", roomID).
		Where("day_of_week = ?", dayOfWeek).
		Where("start_time <= ?", currentTime).
		Where("end_time >= ?", currentTime).
		First(&lesson).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[FindByRoomAndTime] 授業が見つかりませんでした")
			return nil, ErrorRecordNotFound
		}
		log.Printf("[FindByRoomAndTime] エラー: %v", err)
		return nil, err
	}

	log.Printf("[FindByRoomAndTime] 授業を発見: LessonID=%s, Subject=%s", lesson.ID, lesson.SubjectID)
	return &lesson, nil
}
