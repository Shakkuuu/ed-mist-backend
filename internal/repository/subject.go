package repository

import (
	"context"
	"errors"

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

// Create 科目を作成
func (r *SubjectRepository) Create(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

// FindByID IDで科目を取得
func (r *SubjectRepository) FindByID(ctx context.Context, id string) (*model.Subject, error) {
	var subject model.Subject
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&subject).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}
	return &subject, nil
}

// FindByOrgID 組織IDで科目一覧を取得
func (r *SubjectRepository) FindByOrgID(ctx context.Context, orgID string) ([]model.Subject, error) {
	var subjects []model.Subject
	err := r.db.WithContext(ctx).Where("org_id = ?", orgID).Find(&subjects).Error
	return subjects, err
}

// FindByOrgIDAndYear 組織IDと年度で科目一覧を取得
func (r *SubjectRepository) FindByOrgIDAndYear(ctx context.Context, orgID string, year int) ([]model.Subject, error) {
	var subjects []model.Subject
	err := r.db.WithContext(ctx).Where("org_id = ? AND year = ?", orgID, year).Find(&subjects).Error
	return subjects, err
}

// FindAll 全科目を取得
func (r *SubjectRepository) FindAll(ctx context.Context) ([]model.Subject, error) {
	var subjects []model.Subject
	err := r.db.WithContext(ctx).Find(&subjects).Error
	return subjects, err
}

// Update 科目を更新
func (r *SubjectRepository) Update(ctx context.Context, subject *model.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

// Delete 科目を削除
func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Subject{}, "id = ?", id).Error
}
