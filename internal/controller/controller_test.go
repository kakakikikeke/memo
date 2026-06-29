package controller_test

import (
	stdctx "context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	beegocontext "github.com/beego/beego/v2/server/web/context"
	basectrl "github.com/kakakikikeke/memo/internal/controller"
	textctrl "github.com/kakakikikeke/memo/internal/controller/text"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
	"github.com/stretchr/testify/assert"
)

type mockSession struct {
	store map[interface{}]interface{}
}

func (m *mockSession) Set(ctx stdctx.Context, key interface{}, val interface{}) error {
	m.store[key] = val
	return nil
}
func (m *mockSession) Get(ctx stdctx.Context, key interface{}) interface{} {
	return m.store[key]
}
func (m *mockSession) Delete(ctx stdctx.Context, key interface{}) error {
	delete(m.store, key)
	return nil
}
func (m *mockSession) SessionID(ctx stdctx.Context) string                               { return "mock-session-id" }
func (m *mockSession) Release(w http.ResponseWriter)                                     {}
func (m *mockSession) Flush(ctx stdctx.Context) error                                    { return nil }
func (m *mockSession) SessionRelease(ctx stdctx.Context, w http.ResponseWriter)          {}
func (m *mockSession) SessionReleaseIfPresent(ctx stdctx.Context, w http.ResponseWriter) {}

type mockRedisClient struct {
	data map[string][]string
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{data: make(map[string][]string)}
}

type mockRepository struct {
	redis *mockRedisClient
}

func (m *mockRepository) List(ctx stdctx.Context, key string) ([]string, error) {
	return m.redis.data[key], nil
}

func (m *mockRepository) Save(ctx stdctx.Context, key string, value string) error {
	m.redis.data[key] = append([]string{value}, m.redis.data[key]...)
	return nil
}

func (m *mockRepository) Delete(ctx stdctx.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.redis.data, key)
	}
	return nil
}

func (m *mockRepository) Get(ctx stdctx.Context, key string) (string, error) {
	return "", errors.New("not found")
}

func (m *mockRepository) Set(ctx stdctx.Context, key string, value string, expiration time.Duration) error {
	return nil
}

func (m *mockRepository) Text() repository.TextRepository {
	return m
}

func (m *mockRepository) Image() repository.ImageRepository {
	return m
}

func (m *mockRepository) File() repository.FileRepository {
	return m
}

func (m *mockRepository) User() repository.UserRepository {
	return m
}

func TestSaveText(t *testing.T) {
	mockRedis := newMockRedisClient()
	ctrl := &textctrl.Controller{}
	ctrl = textctrl.NewController(service.NewMemoService(&mockRepository{redis: mockRedis}))

	form := url.Values{}
	form.Add("msg", "Hello, world!")

	req := httptest.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := beegocontext.NewContext()
	ctx.Reset(w, req)

	sess := &mockSession{store: make(map[interface{}]interface{})}
	sess.Set(stdctx.Background(), "user", "testuser")
	ctx.Input.CruSession = sess

	ctrl.Init(ctx, "TextController", "Save", nil)
	ctrl.EnableRender = false
	ctrl.Save()

	assert.Equal(t, []string{"Hello, world!"}, mockRedis.data["testuser:"+"memo"])
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestIsValidFileDataURL(t *testing.T) {
	valid := "data:text/plain;base64,SGVsbG8="
	invalidScheme := "javascript:alert(1)"
	invalidMIME := "data:text/html;base64,PHNjcmlwdD4="
	invalidBase64 := "data:text/plain;base64,%%%"

	assert.True(t, basectrl.IsValidFileDataURL(valid))
	assert.False(t, basectrl.IsValidFileDataURL(invalidScheme))
	assert.False(t, basectrl.IsValidFileDataURL(invalidMIME))
	assert.False(t, basectrl.IsValidFileDataURL(invalidBase64))
}

func TestGetContentReturnsOnlyValidatedURL(t *testing.T) {
	validFileInfo := "data:text/plain;base64,SGVsbG8=^_^hello.txt"
	invalidFileInfo := "javascript:alert(1)^_^bad.txt"

	assert.Equal(t, "data:text/plain;base64,SGVsbG8=", basectrl.GetContent(validFileInfo))
	assert.Equal(t, "", basectrl.GetContent(invalidFileInfo))
}

func TestTextTemplateDoesNotUseStr2HTML(t *testing.T) {
	tplPath := filepath.Join("..", "..", "views", "text.tpl")
	body, err := os.ReadFile(tplPath)
	assert.NoError(t, err)
	assert.NotContains(t, string(body), "str2html")
}

func TestScriptsTemplateInjectsCSRFHeader(t *testing.T) {
	tplPath := filepath.Join("..", "..", "views", "meta", "scripts.tpl")
	body, err := os.ReadFile(tplPath)
	assert.NoError(t, err)
	content := string(body)
	assert.Contains(t, content, "csrfToken")
	assert.Contains(t, content, "X-Xsrftoken")
}
