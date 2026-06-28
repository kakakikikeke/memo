package controller

import (
	_ctx "context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/kakakikikeke/memo/internal/service"
	"github.com/stretchr/testify/assert"
)

type mockSession struct {
	store map[interface{}]interface{}
}

func (m *mockSession) Set(ctx _ctx.Context, key interface{}, val interface{}) error {
	m.store[key] = val
	return nil
}
func (m *mockSession) Get(ctx _ctx.Context, key interface{}) interface{} {
	return m.store[key]
}
func (m *mockSession) Delete(ctx _ctx.Context, key interface{}) error {
	delete(m.store, key)
	return nil
}
func (m *mockSession) SessionID(ctx _ctx.Context) string                               { return "mock-session-id" }
func (m *mockSession) Release(w http.ResponseWriter)                                   {}
func (m *mockSession) Flush(ctx _ctx.Context) error                                    { return nil }
func (m *mockSession) SessionRelease(ctx _ctx.Context, w http.ResponseWriter)          {}
func (m *mockSession) SessionReleaseIfPresent(ctx _ctx.Context, w http.ResponseWriter) {}

func TestSaveText(t *testing.T) {
	mockRedis := NewMockRedisClient()
	ctrl := &MainController{}
	ctrl.SetMemoService(service.NewMemoService(&mockRepository{redis: mockRedis}))

	form := url.Values{}
	form.Add("msg", "Hello, world!")

	req := httptest.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctx := context.NewContext()
	ctx.Reset(w, req)

	sess := &mockSession{store: make(map[interface{}]interface{})}
	sess.Set(_ctx.Background(), "user", "testuser")
	ctx.Input.CruSession = sess

	ctrl.Init(ctx, "MainController", "SaveText", nil)
	ctrl.EnableRender = false
	ctrl.SaveText()

	assert.Equal(t, []string{"Hello, world!"}, mockRedis.data["testuser:"+KeyMemo])
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}
