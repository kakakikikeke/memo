package main

import (
	_ctx "context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

// モックセッションストア
type MockSession struct {
	store map[interface{}]interface{}
}

func (m *MockSession) Set(ctx _ctx.Context, key interface{}, val interface{}) error {
	m.store[key] = val
	return nil
}
func (m *MockSession) Get(ctx _ctx.Context, key interface{}) interface{} {
	return m.store[key]
}
func (m *MockSession) Delete(ctx _ctx.Context, key interface{}) error {
	delete(m.store, key)
	return nil
}
func (m *MockSession) SessionID(ctx _ctx.Context) string                               { return "mock-session-id" }
func (m *MockSession) Release(w http.ResponseWriter)                                   {}
func (m *MockSession) Flush(ctx _ctx.Context) error                                    { return nil }
func (m *MockSession) SessionRelease(ctx _ctx.Context, w http.ResponseWriter)          {}
func (m *MockSession) SessionReleaseIfPresent(ctx _ctx.Context, w http.ResponseWriter) {}

func TestSaveText(t *testing.T) {
	// Replace RedisClient with a mock
	mockRedis := NewMockRedisClient()
	RedisClientFactory = func() RedisClient {
		return mockRedis
	}

	// Create request form data
	form := url.Values{}
	form.Add("msg", "Hello, world!")

	// Create a request and response recorder
	req := httptest.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Prepare Beego context
	ctx := context.NewContext()
	ctx.Reset(w, req)

	// Inject a mock session
	mockSess := &MockSession{store: make(map[interface{}]interface{})}
	mockSess.Set(_ctx.Background(), "user", "testuser")
	ctx.Input.CruSession = mockSess

	// Set the context in the Controller
	controller := &MainController{}
	controller.Init(ctx, "MainController", "SaveText", nil)
	controller.EnableRender = false

	// Execute the function under test
	controller.SaveText()

	// Check if data has been added to Redis
	key := "testuser:" + KEY
	assert.Equal(t, []string{"Hello, world!"}, mockRedis.data[key])

	// Check if the redirect was successful
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}
