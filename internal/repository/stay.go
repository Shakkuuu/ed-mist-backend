package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// StayRepository 滞在リポジトリ
type StayRepository struct {
	db *gorm.DB
}

// NewStayRepository 滞在リポジトリを作成
func NewStayRepository(db *gorm.DB) *StayRepository {
	return &StayRepository{db: db}
}

// Create 滞在を作成
func (r *StayRepository) Create(ctx context.Context, stay *model.Stay) error {
	return r.db.WithContext(ctx).Create(stay).Error
}

// FindByID IDで滞在を取得
func (r *StayRepository) FindByID(ctx context.Context, id int) (*model.Stay, error) {
	var stay model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail", "created_at", "updated_at", "deleted_at")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "caption", "mist_zone_id", "created_at", "updated_at")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year", "org_id", "created_at", "updated_at")
		}).
		Preload("Lesson", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "subject_id", "room_id", "org_id", "day_of_week", "start_time", "end_time", "period", "created_at", "updated_at")
		}).
		Preload("Lesson.Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year", "org_id")
		}).
		Preload("Lesson.Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("id = ?", id).
		First(&stay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &stay, nil
}

// FindByUserID ユーザーIDで滞在一覧を取得
func (r *StayRepository) FindByUserID(ctx context.Context, userID string) ([]model.Stay, error) {
	var stays []model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "mist_zone_id")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Lesson", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "subject_id", "room_id", "day_of_week", "start_time", "end_time", "period")
		}).
		Where("user_id = ?", userID).
		Find(&stays).Error
	return stays, err
}

// FindActiveByUserID ユーザーIDでアクティブな滞在を取得
func (r *StayRepository) FindActiveByUserID(ctx context.Context, userID string) (*model.Stay, error) {
	var stay model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "caption", "mist_zone_id")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year", "org_id")
		}).
		Preload("Lesson", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "subject_id", "room_id", "org_id", "day_of_week", "start_time", "end_time", "period")
		}).
		Preload("Lesson.Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Lesson.Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "mist_zone_id")
		}).
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&stay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &stay, nil
}

// FindActiveByUserAndLesson ユーザーIDとLessonIDでアクティブな滞在を取得
func (r *StayRepository) FindActiveByUserAndLesson(ctx context.Context, userID string, lessonID string) (*model.Stay, error) {
	var stay model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("user_id = ? AND lesson_id = ? AND is_active = ?", userID, lessonID, true).
		First(&stay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &stay, nil
}

// FindByLessonID LessonIDで滞在一覧を取得
func (r *StayRepository) FindByLessonID(ctx context.Context, lessonID string) ([]model.Stay, error) {
	var stays []model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "org_room_id", "name", "mist_zone_id")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "year")
		}).
		Preload("Lesson", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "subject_id", "room_id", "day_of_week", "start_time", "end_time", "period")
		}).
		Where("lesson_id = ?", lessonID).
		Find(&stays).Error
	return stays, err
}

// FindByRoomID 部屋IDで滞在一覧を取得
func (r *StayRepository) FindByRoomID(ctx context.Context, roomID string) ([]model.Stay, error) {
	var stays []model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "org_id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "mist_zone_id")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Lesson", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "start_time", "end_time", "period")
		}).
		Where("room_id = ?", roomID).
		Find(&stays).Error
	return stays, err
}

// FindActiveByRoomID 部屋IDでアクティブな滞在一覧を取得
func (r *StayRepository) FindActiveByRoomID(ctx context.Context, roomID string) ([]model.Stay, error) {
	var stays []model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("room_id = ? AND is_active = ?", roomID, true).
		Find(&stays).Error
	return stays, err
}

// FindAll 全滞在を取得
func (r *StayRepository) FindAll(ctx context.Context) ([]model.Stay, error) {
	var stays []model.Stay
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mail")
		}).
		Preload("Room", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Find(&stays).Error
	return stays, err
}

// Update 滞在を更新
func (r *StayRepository) Update(ctx context.Context, stay *model.Stay) error {
	return r.db.WithContext(ctx).Save(stay).Error
}

// Delete 滞在を削除
func (r *StayRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Stay{}, "id = ?", id).Error
}

// EndStay 滞在を終了する
func (r *StayRepository) EndStay(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Model(&model.Stay{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": false,
		"leaved_at": "NOW()",
	}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorRecordNotFound
		}
		return err
	}
	return nil
}
