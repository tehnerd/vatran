package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func setupAuthenticatorTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "authenticator_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("NewStore failed: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.RemoveAll(tmpDir)
	}
	return store, cleanup
}

func TestNewBasicAuthenticator(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, true, []string{"/api/public/*"}, true)

	if auth == nil {
		t.Fatal("NewBasicAuthenticator returned nil")
	}

	// Check default exempt paths are included
	exemptPaths := auth.exemptPaths
	hasLogin := false
	hasHealth := false
	hasPublic := false
	for _, p := range exemptPaths {
		if p == "/login" {
			hasLogin = true
		}
		if p == "/health" {
			hasHealth = true
		}
		if p == "/api/public/*" {
			hasPublic = true
		}
	}

	if !hasLogin {
		t.Error("/login not in exempt paths")
	}
	if !hasHealth {
		t.Error("/health not in exempt paths")
	}
	if !hasPublic {
		t.Error("/api/public/* not in exempt paths")
	}
}

func TestBasicAuthenticator_Authenticate(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, true, []string{"/static/*"}, false)

	// Create a test user and session
	user, err := auth.CreateUser("testuser", "password")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	token, _, err := auth.Login("testuser", "password")
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	t.Run("allows exempt path /login", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("exempt path /login was not allowed")
		}
		if result.UserID != "anonymous" {
			t.Errorf("UserID = %q, want %q", result.UserID, "anonymous")
		}
	})

	t.Run("allows exempt path /health", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("exempt path /health was not allowed")
		}
	})

	t.Run("allows wildcard exempt path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/js/app.js", nil)
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("wildcard exempt path was not allowed")
		}
	})

	t.Run("allows localhost when configured", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("localhost was not allowed")
		}
		if result.UserID != "localhost" {
			t.Errorf("UserID = %q, want %q", result.UserID, "localhost")
		}
	})

	t.Run("allows IPv6 localhost", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "[::1]:12345"
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("IPv6 localhost was not allowed")
		}
	})

	t.Run("denies request without cookie from non-localhost", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		result := auth.Authenticate(req)
		if result.Allowed {
			t.Error("request without cookie was incorrectly allowed")
		}
		if result.Message != "authentication required" {
			t.Errorf("Message = %q, want %q", result.Message, "authentication required")
		}
	})

	t.Run("allows request with valid session cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: token,
		})
		result := auth.Authenticate(req)
		if !result.Allowed {
			t.Error("valid session was not allowed")
		}
		if result.UserID != user.Username {
			t.Errorf("UserID = %q, want %q", result.UserID, user.Username)
		}
	})

	t.Run("denies request with invalid session cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: "invalidtoken",
		})
		result := auth.Authenticate(req)
		if result.Allowed {
			t.Error("invalid session was incorrectly allowed")
		}
		if result.Message != "invalid or expired session" {
			t.Errorf("Message = %q, want %q", result.Message, "invalid or expired session")
		}
	})
}

func TestBasicAuthenticator_Authenticate_LocalhostDisabled(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	t.Run("denies localhost when disabled", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		result := auth.Authenticate(req)
		if result.Allowed {
			t.Error("localhost was allowed when disabled")
		}
	})
}

func TestBasicAuthenticator_Login(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	_, err := auth.CreateUser("loginuser", "correctpass")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Run("successful login", func(t *testing.T) {
		token, expiresAt, err := auth.Login("loginuser", "correctpass")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if token == "" {
			t.Error("token is empty")
		}
		if expiresAt.Before(time.Now()) {
			t.Error("expiration is in the past")
		}
	})

	t.Run("failed login wrong password", func(t *testing.T) {
		_, _, err := auth.Login("loginuser", "wrongpass")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("failed login nonexistent user", func(t *testing.T) {
		_, _, err := auth.Login("nouser", "pass")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}

func TestBasicAuthenticator_Logout(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	_, _ = auth.CreateUser("logoutuser", "pass")
	token, _, _ := auth.Login("logoutuser", "pass")

	err := auth.Logout(token)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Verify session no longer works
	req := httptest.NewRequest("GET", "/protected", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: token,
	})
	result := auth.Authenticate(req)
	if result.Allowed {
		t.Error("logged out session still works")
	}
}

func TestBasicAuthenticator_CreateUser(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	t.Run("creates user successfully", func(t *testing.T) {
		user, err := auth.CreateUser("newuser", "password")
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
		if user.Username != "newuser" {
			t.Errorf("Username = %q, want %q", user.Username, "newuser")
		}
	})

	t.Run("returns error for duplicate user", func(t *testing.T) {
		_, err := auth.CreateUser("dupuser", "pass1")
		if err != nil {
			t.Fatalf("first CreateUser failed: %v", err)
		}

		_, err = auth.CreateUser("dupuser", "pass2")
		if !errors.Is(err, ErrUserExists) {
			t.Errorf("expected ErrUserExists, got: %v", err)
		}
	})
}

func TestBasicAuthenticator_UserCount(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	count, err := auth.UserCount()
	if err != nil {
		t.Fatalf("UserCount failed: %v", err)
	}
	if count != 0 {
		t.Errorf("initial count = %d, want 0", count)
	}

	_, _ = auth.CreateUser("user1", "pass")
	_, _ = auth.CreateUser("user2", "pass")

	count, err = auth.UserCount()
	if err != nil {
		t.Fatalf("UserCount failed: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestBasicAuthenticator_IsTLS(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	authTLS := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, true)
	if !authTLS.IsTLS() {
		t.Error("IsTLS() = false, want true")
	}

	store2, cleanup2 := setupAuthenticatorTestStore(t)
	defer cleanup2()

	authNoTLS := NewBasicAuthenticator(store2, bcrypt.MinCost, 24, false, nil, false)
	if authNoTLS.IsTLS() {
		t.Error("IsTLS() = true, want false")
	}
}

func TestBasicAuthenticator_CleanupExpiredSessions(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 1, false, nil, false)

	_, _ = auth.CreateUser("cleanupuser", "pass")
	token, _, _ := auth.Login("cleanupuser", "pass")

	// Manually expire the session
	_, _ = store.DB().Exec(
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		time.Now().Add(-1*time.Hour), token,
	)

	count, err := auth.CleanupExpiredSessions()
	if err != nil {
		t.Fatalf("CleanupExpiredSessions failed: %v", err)
	}
	if count != 1 {
		t.Errorf("cleaned up count = %d, want 1", count)
	}
}

func TestBasicAuthenticator_SetSessionCookie(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	t.Run("TLS cookie", func(t *testing.T) {
		auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, true)

		recorder := httptest.NewRecorder()
		expiresAt := time.Now().Add(24 * time.Hour)
		auth.SetSessionCookie(recorder, "testtoken", expiresAt)

		cookies := recorder.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("expected 1 cookie, got %d", len(cookies))
		}

		cookie := cookies[0]
		if cookie.Name != SessionCookieName {
			t.Errorf("cookie name = %q, want %q", cookie.Name, SessionCookieName)
		}
		if cookie.Value != "testtoken" {
			t.Errorf("cookie value = %q, want %q", cookie.Value, "testtoken")
		}
		if !cookie.HttpOnly {
			t.Error("cookie is not HttpOnly")
		}
		if !cookie.Secure {
			t.Error("TLS cookie is not Secure")
		}
		if cookie.SameSite != http.SameSiteStrictMode {
			t.Errorf("SameSite = %v, want %v", cookie.SameSite, http.SameSiteStrictMode)
		}
	})

	t.Run("non-TLS cookie", func(t *testing.T) {
		store2, cleanup2 := setupAuthenticatorTestStore(t)
		defer cleanup2()

		auth := NewBasicAuthenticator(store2, bcrypt.MinCost, 24, false, nil, false)

		recorder := httptest.NewRecorder()
		auth.SetSessionCookie(recorder, "testtoken", time.Now().Add(time.Hour))

		cookie := recorder.Result().Cookies()[0]
		if cookie.Secure {
			t.Error("non-TLS cookie should not be Secure")
		}
	})
}

func TestBasicAuthenticator_ClearSessionCookie(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, true)

	recorder := httptest.NewRecorder()
	auth.ClearSessionCookie(recorder)

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != SessionCookieName {
		t.Errorf("cookie name = %q, want %q", cookie.Name, SessionCookieName)
	}
	if cookie.Value != "" {
		t.Errorf("cookie value = %q, want empty", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("MaxAge = %d, want -1", cookie.MaxAge)
	}
}

func TestBasicAuthenticator_isExemptPath(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, []string{"/api/public/*", "/custom"}, false)

	tests := []struct {
		path   string
		exempt bool
	}{
		{"/login", true},
		{"/health", true},
		{"/custom", true},
		{"/api/public/data", true},
		{"/api/public/nested/path", true},
		{"/api/private", false},
		{"/protected", false},
		{"/loginextra", false}, // Should not match /login prefix
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			result := auth.isExemptPath(tc.path)
			if result != tc.exempt {
				t.Errorf("isExemptPath(%q) = %v, want %v", tc.path, result, tc.exempt)
			}
		})
	}
}

func TestBasicAuthenticator_isLocalhost(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, true, nil, false)

	tests := []struct {
		remoteAddr  string
		isLocalhost bool
	}{
		{"127.0.0.1:8080", true},
		{"[::1]:8080", true},
		{"192.168.1.1:8080", false},
		{"10.0.0.1:8080", false},
		{"[2001:db8::1]:8080", false},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.remoteAddr, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tc.remoteAddr
			result := auth.isLocalhost(req)
			if result != tc.isLocalhost {
				t.Errorf("isLocalhost(%q) = %v, want %v", tc.remoteAddr, result, tc.isLocalhost)
			}
		})
	}
}

func TestBasicAuthenticator_Close(t *testing.T) {
	store, cleanup := setupAuthenticatorTestStore(t)
	defer cleanup()

	auth := NewBasicAuthenticator(store, bcrypt.MinCost, 24, false, nil, false)

	err := auth.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
