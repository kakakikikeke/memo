package user_test

import (
	stdctx "context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	beegocontext "github.com/beego/beego/v2/server/web/context"
	userctrl "github.com/kakakikikeke/memo/internal/controller/user"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

type inMemoryRepository struct {
	users map[string]string
	data  map[string][]string
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{users: make(map[string]string), data: make(map[string][]string)}
}

func (r *inMemoryRepository) List(ctx stdctx.Context, key string) ([]string, error) {
	return append([]string(nil), r.data[key]...), nil
}

func (r *inMemoryRepository) Save(ctx stdctx.Context, key string, value string) error {
	r.data[key] = append([]string{value}, r.data[key]...)
	return nil
}

func (r *inMemoryRepository) Delete(ctx stdctx.Context, keys ...string) error {
	for _, key := range keys {
		delete(r.data, key)
		delete(r.users, key)
	}
	return nil
}

func (r *inMemoryRepository) Get(ctx stdctx.Context, key string) (string, error) {
	v, ok := r.users[key]
	if !ok {
		return "", redis.Nil
	}
	return v, nil
}

func (r *inMemoryRepository) Set(ctx stdctx.Context, key string, value string, expiration time.Duration) error {
	r.users[key] = value
	return nil
}

func (r *inMemoryRepository) Text() repository.TextRepository {
	return r
}

func (r *inMemoryRepository) Image() repository.ImageRepository {
	return r
}

func (r *inMemoryRepository) File() repository.FileRepository {
	return r
}

func (r *inMemoryRepository) User() repository.UserRepository {
	return r
}

func TestCreatePasswordMismatchDoesNotCreateUser(t *testing.T) {
	repo := newInMemoryRepository()
	ctrl := userctrl.NewController(service.NewMemoService(repo))

	form := url.Values{}
	form.Add("name", "alice")
	form.Add("pass", "secret")
	form.Add("pass2", "different")

	req := httptest.NewRequest("POST", "/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := beegocontext.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{store: make(map[interface{}]interface{})}

	ctrl.Init(ctx, "UserController", "Create", nil)
	ctrl.EnableRender = false
	ctrl.Create()

	assert.Equal(t, 403, w.Code)
	_, exists := repo.users["alice"]
	assert.False(t, exists)
}

func TestCreateUserSuccessCreatesUserAndSession(t *testing.T) {
	repo := newInMemoryRepository()
	ctrl := userctrl.NewController(service.NewMemoService(repo))

	form := url.Values{}
	form.Add("name", "alice")
	form.Add("pass", "secret")
	form.Add("pass2", "secret")

	req := httptest.NewRequest("POST", "/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := beegocontext.NewContext()
	ctx.Reset(w, req)
	sess := &mockSession{store: make(map[interface{}]interface{})}
	ctx.Input.CruSession = sess

	ctrl.Init(ctx, "UserController", "Create", nil)
	ctrl.EnableRender = false
	ctrl.Create()

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "alice", sess.Get(stdctx.Background(), "user"))
	hash, exists := repo.users["alice"]
	assert.True(t, exists)
	assert.NotEqual(t, "secret", hash)
}

func TestCheckWrongPasswordReturnsForbidden(t *testing.T) {
	repo := newInMemoryRepository()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	repo.users["alice"] = string(hash)

	ctrl := userctrl.NewController(service.NewMemoService(repo))

	form := url.Values{}
	form.Add("name", "alice")
	form.Add("pass", "wrong")

	req := httptest.NewRequest("POST", "/check", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := beegocontext.NewContext()
	ctx.Reset(w, req)
	sess := &mockSession{store: make(map[interface{}]interface{})}
	ctx.Input.CruSession = sess

	ctrl.Init(ctx, "UserController", "Check", nil)
	ctrl.EnableRender = false
	ctrl.Check()

	assert.Equal(t, 403, w.Code)
	assert.Nil(t, sess.Get(stdctx.Background(), "user"))
}

func TestCreateInvalidFormReturnsBadRequest(t *testing.T) {
	repo := newInMemoryRepository()
	ctrl := userctrl.NewController(service.NewMemoService(repo))

	req := httptest.NewRequest("POST", "/create", strings.NewReader("name=%zz&pass=x&pass2=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := beegocontext.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{store: make(map[interface{}]interface{})}

	ctrl.Init(ctx, "UserController", "Create", nil)
	ctrl.EnableRender = false
	ctrl.Create()

	assert.Equal(t, 400, w.Code)
}
