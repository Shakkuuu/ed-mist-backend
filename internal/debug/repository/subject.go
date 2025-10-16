package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

// SubjectRepository 科目リポジトリ
type SubjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository 科目リポジトリを作成
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// FindAll 全科目を取得
func (r *SubjectRepository) FindAll(ctx context.Context) ([]model.Subject, error) {
	var subjects []model.Subject
	err := r.db.WithContext(ctx).Find(&subjects).Error
	return subjects, err
}

// Create 科目を作成
func (r *SubjectRepository) Create(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

// DeleteAll 全科目を削除
func (r *SubjectRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM subjects").Error
}
