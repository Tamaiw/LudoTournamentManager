package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	registerFn     func(ctx context.Context, email, password, inviteCode string) (string, error)
	loginFn        func(ctx context.Context, email, password string) (string, error)
	logoutFn       func(ctx context.Context, token string) error
	getCurrentFn   func(ctx context.Context, token string) (*inbound.UserDTO, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password, inviteCode string) (string, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, email, password, inviteCode)
	}
	return "mock-token", nil
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, email, password)
	}
	return "mock-token", nil
}

func (m *mockAuthService) Logout(ctx context.Context, token string) error {
	if m.logoutFn != nil {
		return m.logoutFn(ctx, token)
	}
	return nil
}

func (m *mockAuthService) GetCurrentUser(ctx context.Context, token string) (*inbound.UserDTO, error) {
	if m.getCurrentFn != nil {
		return m.getCurrentFn(ctx, token)
	}
	return &inbound.UserDTO{ID: "user-1", Email: "test@example.com", Role: "member"}, nil
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterHandler_Success(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(ctx context.Context, email, password, inviteCode string) (string, error) {
			return "new-token", nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"email":"test@example.com","password":"password123","inviteCode":"CODE123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := RegisterHandler(svc)
	handler(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp["token"] != "new-token" {
		t.Errorf("expected token 'new-token', got %q", resp["token"])
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	svc := &mockAuthService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := RegisterHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	svc := &mockAuthService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Missing password
	body := `{"email":"test@example.com"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := RegisterHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	svc := &mockAuthService{
		loginFn: func(ctx context.Context, email, password string) (string, error) {
			return "login-token", nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"email":"test@example.com","password":"password123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := LoginHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp["token"] != "login-token" {
		t.Errorf("expected token 'login-token', got %q", resp["token"])
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	svc := &mockAuthService{
		loginFn: func(ctx context.Context, email, password string) (string, error) {
			return "", gin.Error{Err: nil}
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"email":"test@example.com","password":"wrongpassword"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := LoginHandler(svc)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLogoutHandler_MissingToken(t *testing.T) {
	svc := &mockAuthService{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/auth/logout", nil)

	handler := LogoutHandler(svc)
	handler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMeHandler_Success(t *testing.T) {
	svc := &mockAuthService{
		getCurrentFn: func(ctx context.Context, token string) (*inbound.UserDTO, error) {
			return &inbound.UserDTO{ID: "user-123", Email: "user@example.com", Role: "admin"}, nil
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/auth/me", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	handler := MeHandler(svc)
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var user inbound.UserDTO
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}
	if user.ID != "user-123" {
		t.Errorf("expected user ID 'user-123', got %q", user.ID)
	}
}

func TestMeHandler_Unauthorized(t *testing.T) {
	svc := &mockAuthService{
		getCurrentFn: func(ctx context.Context, token string) (*inbound.UserDTO, error) {
			return nil, gin.Error{Err: nil}
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/auth/me", nil)

	handler := MeHandler(svc)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}