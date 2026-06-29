package user

import (
	"context"

	"github.com/kakakikikeke/memo/internal/model"
	"github.com/kakakikikeke/memo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (bool, error) {
	hash, err := s.repo.Get(ctx, username)
	if err != nil {
		return false, err
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}

func (s *Service) Create(ctx context.Context, username, password string) error {
	if _, err := s.repo.Get(ctx, username); err == nil {
		return ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Set(ctx, username, string(hash), 0)
}

func (s *Service) Delete(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyMemo), model.UserKey(username, model.KeyFile), model.UserKey(username, model.KeyImage), username)
}
