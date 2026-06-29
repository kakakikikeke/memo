package service

import (
	"context"
	"strings"

	"github.com/kakakikikeke/memo/internal/model"
	"github.com/kakakikikeke/memo/internal/repository"
	filesvc "github.com/kakakikikeke/memo/internal/service/file"
	imagesvc "github.com/kakakikikeke/memo/internal/service/image"
	textsvc "github.com/kakakikikeke/memo/internal/service/text"
	usersvc "github.com/kakakikikeke/memo/internal/service/user"
)

var (
	ErrTextLimitExceeded  = textsvc.ErrTextLimitExceeded
	ErrImageLimitExceeded = imagesvc.ErrImageLimitExceeded
	ErrFileLimitExceeded  = filesvc.ErrFileLimitExceeded
	ErrUserAlreadyExists  = usersvc.ErrUserAlreadyExists
)

type MemoService struct {
	repo         repository.Provider
	textService  *textsvc.Service
	imageService *imagesvc.Service
	fileService  *filesvc.Service
	userService  *usersvc.Service
}

func NewMemoService(repo repository.Provider) *MemoService {
	return &MemoService{repo: repo}
}

func (s *MemoService) textSvc() *textsvc.Service {
	if s.textService == nil {
		s.textService = textsvc.NewService(s.repo.Text())
	}
	return s.textService
}

func (s *MemoService) imageSvc() *imagesvc.Service {
	if s.imageService == nil {
		s.imageService = imagesvc.NewService(s.repo.Image())
	}
	return s.imageService
}

func (s *MemoService) fileSvc() *filesvc.Service {
	if s.fileService == nil {
		s.fileService = filesvc.NewService(s.repo.File())
	}
	return s.fileService
}

func (s *MemoService) userSvc() *usersvc.Service {
	if s.userService == nil {
		s.userService = usersvc.NewService(s.repo.User())
	}
	return s.userService
}

func (s *MemoService) ListText(ctx context.Context, username string) ([]string, error) {
	return s.textSvc().List(ctx, username)
}

func (s *MemoService) SaveText(ctx context.Context, username string, msg string) error {
	return s.textSvc().Save(ctx, username, msg)
}

func (s *MemoService) ClearText(ctx context.Context, username string) error {
	return s.textSvc().Clear(ctx, username)
}

func (s *MemoService) ListImage(ctx context.Context, username string) ([]string, error) {
	return s.imageSvc().List(ctx, username)
}

func (s *MemoService) SaveImage(ctx context.Context, username string, base64 string) error {
	return s.imageSvc().Save(ctx, username, base64)
}

func (s *MemoService) ClearImage(ctx context.Context, username string) error {
	return s.imageSvc().Clear(ctx, username)
}

func (s *MemoService) ListFile(ctx context.Context, username string) ([]string, error) {
	return s.fileSvc().List(ctx, username)
}

func (s *MemoService) SaveFile(ctx context.Context, username string, value string) error {
	return s.fileSvc().Save(ctx, username, value)
}

func (s *MemoService) ClearFile(ctx context.Context, username string) error {
	return s.fileSvc().Clear(ctx, username)
}

func (s *MemoService) Authenticate(ctx context.Context, username, password string) (bool, error) {
	return s.userSvc().Authenticate(ctx, username, password)
}

func (s *MemoService) CreateUser(ctx context.Context, username, password string) error {
	return s.userSvc().Create(ctx, username, password)
}

func (s *MemoService) DeleteUser(ctx context.Context, username string) error {
	return s.userSvc().Delete(ctx, username)
}

func (s *MemoService) SaveTextWithLimit(ctx context.Context, username, msg string) error {
	return s.textSvc().SaveWithLimit(ctx, username, msg)
}

func (s *MemoService) SaveImageWithLimit(ctx context.Context, username, base64 string) error {
	return s.imageSvc().SaveWithLimit(ctx, username, base64)
}

func (s *MemoService) SaveFileWithLimit(ctx context.Context, username, value string) error {
	return s.fileSvc().SaveWithLimit(ctx, username, value)
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
