package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kakakikikeke/memo/internal/model"
	"github.com/kakakikikeke/memo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTextLimitExceeded  = errors.New("exceeds the number of texts that can be uploaded")
	ErrImageLimitExceeded = errors.New("exceeds the number of images that can be uploaded")
	ErrFileLimitExceeded  = errors.New("exceeds the number of files that can be uploaded")
	ErrUserAlreadyExists  = errors.New("specified user already exists")
)

type MemoService struct {
	repo repository.Repository
}

func NewMemoService(repo repository.Repository) *MemoService {
	return &MemoService{repo: repo}
}

func (s *MemoService) ListText(ctx context.Context, username string) ([]string, error) {
	return s.repo.List(ctx, model.UserKey(username, model.KeyMemo))
}

func (s *MemoService) SaveText(ctx context.Context, username string, msg string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyMemo), msg)
}

func (s *MemoService) ClearText(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyMemo))
}

func (s *MemoService) ListImage(ctx context.Context, username string) ([]string, error) {
	return s.repo.List(ctx, model.UserKey(username, model.KeyImage))
}

func (s *MemoService) SaveImage(ctx context.Context, username string, base64 string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyImage), base64)
}

func (s *MemoService) ClearImage(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyImage))
}

func (s *MemoService) ListFile(ctx context.Context, username string) ([]string, error) {
	return s.repo.List(ctx, model.UserKey(username, model.KeyFile))
}

func (s *MemoService) SaveFile(ctx context.Context, username string, value string) error {
	return s.repo.Save(ctx, model.UserKey(username, model.KeyFile), value)
}

func (s *MemoService) ClearFile(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyFile))
}

func (s *MemoService) Authenticate(ctx context.Context, username, password string) (bool, error) {
	hash, err := s.repo.Get(ctx, username)
	if err != nil {
		return false, err
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}

func (s *MemoService) CreateUser(ctx context.Context, username, password string) error {
	if _, err := s.repo.Get(ctx, username); err == nil {
		return ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Set(ctx, username, string(hash), 0)
}

func (s *MemoService) DeleteUser(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, model.UserKey(username, model.KeyMemo), model.UserKey(username, model.KeyFile), model.UserKey(username, model.KeyImage), username)
}

func (s *MemoService) SaveTextWithLimit(ctx context.Context, username, msg string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyMemo))
	if err != nil {
		return err
	}
	if len(items) >= 10 {
		return ErrTextLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyMemo), msg)
}

func (s *MemoService) SaveImageWithLimit(ctx context.Context, username, base64 string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyImage))
	if err != nil {
		return err
	}
	if len(items) >= 1 {
		return ErrImageLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyImage), base64)
}

func (s *MemoService) SaveFileWithLimit(ctx context.Context, username, value string) error {
	items, err := s.repo.List(ctx, model.UserKey(username, model.KeyFile))
	if err != nil {
		return err
	}
	if len(items) >= 1 {
		return ErrFileLimitExceeded
	}
	return s.repo.Save(ctx, model.UserKey(username, model.KeyFile), value)
}

func (s *MemoService) ReplaceSpaceWithPlus(value string) string {
	return strings.ReplaceAll(value, " ", "+")
}

func (s *MemoService) FormatFileValue(base64, filename string) string {
	return s.ReplaceSpaceWithPlus(base64) + "^_^" + filename
}

func (s *MemoService) UserKey(username string, kind string) string {
	return model.UserKey(username, kind)
}
