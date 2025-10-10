package service

import (
	"context"
	"time"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
	"github.com/Shakkuuu/ed-mist-backend/internal/repository"

	"github.com/google/uuid"
)

// SubjectService 科目サービス
type SubjectService struct {
	subjectRepo *repository.SubjectRepository
}

// NewSubjectService 科目サービスを作成
func NewSubjectService(subjectRepo *repository.SubjectRepository) *SubjectService {
	return &SubjectService{
		subjectRepo: subjectRepo,
	}
}

// Create 科目を作成
func (s *SubjectService) Create(ctx context.Context, name string, year int, orgID string) (*model.Subject, error) {
	subject := &model.Subject{
		ID:        uuid.NewString(),
		Name:      name,
		Year:      year,
		OrgID:     orgID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

// GetAll 全科目を取得
func (s *SubjectService) GetAll(ctx context.Context) ([]model.Subject, error) {
	subjects, err := s.subjectRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

// GetByID IDで科目を取得
func (s *SubjectService) GetByID(ctx context.Context, id string) (*model.Subject, error) {
	subject, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

// GetByOrgID 組織IDで科目一覧を取得
func (s *SubjectService) GetByOrgID(ctx context.Context, orgID string) ([]model.Subject, error) {
	subjects, err := s.subjectRepo.FindByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

// GetByOrgIDAndYear 組織IDと年度で科目一覧を取得
func (s *SubjectService) GetByOrgIDAndYear(ctx context.Context, orgID string, year int) ([]model.Subject, error) {
	subjects, err := s.subjectRepo.FindByOrgIDAndYear(ctx, orgID, year)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

// Update 科目を更新
func (s *SubjectService) Update(ctx context.Context, subject *model.Subject, name string, year int) error {
	subject.Name = name
	subject.Year = year
	subject.UpdatedAt = time.Now()

	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		return err
	}
	return nil
}

// Delete 科目を削除
func (s *SubjectService) Delete(ctx context.Context, id string) error {
	if err := s.subjectRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
