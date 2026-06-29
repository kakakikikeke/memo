package file
package file

import (
	"context"

	"github.com/kakakikikeke/memo/internal/model"
	"github.com/kakakikikeke/memo/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, username string) ([]string, error) {
	return s.repo.List(ctx, model.UserKey(username, model.KeyFile))
}

func (s *Service) Save(ctx context.Context, username string, value string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyFile), value)
}

func (s *Service) Clear(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyFile))
}

func (s *Service) SaveWithLimit(ctx context.Context, username, value string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyFile))
	if err != nil {
		return err
	}
	if len(items) >= 1 {
		return ErrFileLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyFile), value)
}
