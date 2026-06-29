package text
package text

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
	return s.repo.List(ctx, model.UserKey(username, model.KeyMemo))
}

func (s *Service) Save(ctx context.Context, username string, msg string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyMemo), msg)
}

func (s *Service) Clear(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyMemo))
}

func (s *Service) SaveWithLimit(ctx context.Context, username, msg string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyMemo))
	if err != nil {
		return err
	}
	if len(items) >= 10 {
		return ErrTextLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyMemo), msg)
}
