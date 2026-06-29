package image

import (
	"context"

	"github.com/kakakikikeke/memo/internal/model"
	"github.com/kakakikikeke/memo/internal/repository"
)

type Service struct {
	repo repository.ImageRepository
}

func NewService(repo repository.ImageRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, username string) ([]string, error) {
	return s.repo.List(ctx, model.UserKey(username, model.KeyImage))
}

func (s *Service) Save(ctx context.Context, username string, base64 string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyImage), base64)
}

func (s *Service) Clear(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyImage))
}

func (s *Service) SaveWithLimit(ctx context.Context, username, base64 string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyImage))
	if err != nil {
		return err
	}
	if len(items) >= 1 {
		return ErrImageLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyImage), base64)
}
