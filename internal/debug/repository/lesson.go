package repository

import (
	"context"

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

// FindAll 全授業を取得
func (r *LessonRepository) FindAll(ctx context.Context) ([]model.Lesson, error) {
	var lessons []model.Lesson
	err := r.db.WithContext(ctx).Find(&lessons).Error
	return lessons, err
}

// Create 授業を作成
func (r *LessonRepository) Create(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

// DeleteAll 全授業を削除
func (r *LessonRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM lessons").Error
}
