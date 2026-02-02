package auth

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func setupSessionTestStore(t *testing.T) (*Store, int64, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "session_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("NewStore failed: %v", err)
	}

	// Create a test user for sessions
	userRepo := NewUserRepository(store, bcrypt.MinCost)
	user, err := userRepo.Create("sessiontestuser", "password")
	if err != nil {
		store.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to create test user: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.RemoveAll(tmpDir)
	}
	return store, user.ID, cleanup
}

func TestNewSessionRepository(t *testing.T) {
	store, _, cleanup := setupSessionTestStore(t)
	defer cleanup()

	t.Run("valid timeout", func(t *testing.T) {
		repo := NewSessionRepository(store, 24)
		if repo.sessionTimeout != 24*time.Hour {
			t.Errorf("sessionTimeout = %v, want %v", repo.sessionTimeout, 24*time.Hour)
		}
	})

	t.Run("zero timeout defaults to 24 hours", func(t *testing.T) {
		repo := NewSessionRepository(store, 0)
		if repo.sessionTimeout != 24*time.Hour {
			t.Errorf("sessionTimeout = %v, want %v", repo.sessionTimeout, 24*time.Hour)
		}
	})

	t.Run("negative timeout defaults to 24 hours", func(t *testing.T) {
		repo := NewSessionRepository(store, -5)
		if repo.sessionTimeout != 24*time.Hour {
			t.Errorf("sessionTimeout = %v, want %v", repo.sessionTimeout, 24*time.Hour)
		}
	})
}

func TestSessionRepository_Create(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 24)

	t.Run("creates session successfully", func(t *testing.T) {
		session, err := repo.Create(userID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if session.ID == "" {
			t.Error("session ID is empty")
		}
		if len(session.ID) != 64 { // 256-bit = 32 bytes = 64 hex chars
			t.Errorf("session ID length = %d, want 64", len(session.ID))
		}
		if session.UserID != userID {
			t.Errorf("UserID = %d, want %d", session.UserID, userID)
		}
		if session.ExpiresAt.Before(time.Now()) {
			t.Error("session already expired")
		}
		if session.ExpiresAt.After(time.Now().Add(25 * time.Hour)) {
			t.Error("session expiration too far in future")
		}
	})

	t.Run("creates unique tokens", func(t *testing.T) {
		session1, _ := repo.Create(userID)
		session2, _ := repo.Create(userID)
		if session1.ID == session2.ID {
			t.Error("session tokens are not unique")
		}
	})
}

func TestSessionRepository_Get(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 1) // 1 hour timeout

	t.Run("retrieves valid session", func(t *testing.T) {
		created, err := repo.Create(userID)
		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}

		retrieved, err := repo.Get(created.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.ID != created.ID {
			t.Errorf("ID = %q, want %q", retrieved.ID, created.ID)
		}
		if retrieved.UserID != userID {
			t.Errorf("UserID = %d, want %d", retrieved.UserID, userID)
		}
	})

	t.Run("returns ErrSessionNotFound for nonexistent token", func(t *testing.T) {
		_, err := repo.Get("nonexistenttoken")
		if !errors.Is(err, ErrSessionNotFound) {
			t.Errorf("expected ErrSessionNotFound, got: %v", err)
		}
	})

	t.Run("returns ErrSessionExpired for expired session", func(t *testing.T) {
		// Create a session with very short timeout
		shortRepo := NewSessionRepository(store, 0) // defaults to 24h, so we'll manually expire it

		session, _ := shortRepo.Create(userID)

		// Manually expire the session in the database
		_, err := store.DB().Exec(
			"UPDATE sessions SET expires_at = ? WHERE id = ?",
			time.Now().Add(-1*time.Hour), session.ID,
		)
		if err != nil {
			t.Fatalf("failed to expire session: %v", err)
		}

		_, err = shortRepo.Get(session.ID)
		if !errors.Is(err, ErrSessionExpired) {
			t.Errorf("expected ErrSessionExpired, got: %v", err)
		}

		// Verify the session was deleted
		var count int
		store.DB().QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", session.ID).Scan(&count)
		if count != 0 {
			t.Error("expired session was not deleted")
		}
	})
}

func TestSessionRepository_Delete(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 24)

	t.Run("deletes session successfully", func(t *testing.T) {
		session, _ := repo.Create(userID)

		err := repo.Delete(session.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify session is gone
		_, err = repo.Get(session.ID)
		if !errors.Is(err, ErrSessionNotFound) {
			t.Error("session still exists after deletion")
		}
	})

	t.Run("delete nonexistent session does not error", func(t *testing.T) {
		err := repo.Delete("nonexistent")
		if err != nil {
			t.Errorf("Delete returned error for nonexistent session: %v", err)
		}
	})
}

func TestSessionRepository_DeleteByUserID(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 24)

	// Create multiple sessions for the user
	session1, _ := repo.Create(userID)
	session2, _ := repo.Create(userID)
	session3, _ := repo.Create(userID)

	err := repo.DeleteByUserID(userID)
	if err != nil {
		t.Fatalf("DeleteByUserID failed: %v", err)
	}

	// Verify all sessions are gone
	for _, token := range []string{session1.ID, session2.ID, session3.ID} {
		_, err = repo.Get(token)
		if !errors.Is(err, ErrSessionNotFound) {
			t.Errorf("session %s still exists", token)
		}
	}
}

func TestSessionRepository_DeleteExpired(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 24)

	// Create sessions
	validSession, _ := repo.Create(userID)
	expiredSession1, _ := repo.Create(userID)
	expiredSession2, _ := repo.Create(userID)

	// Expire two sessions
	_, _ = store.DB().Exec(
		"UPDATE sessions SET expires_at = ? WHERE id IN (?, ?)",
		time.Now().Add(-1*time.Hour), expiredSession1.ID, expiredSession2.ID,
	)

	count, err := repo.DeleteExpired()
	if err != nil {
		t.Fatalf("DeleteExpired failed: %v", err)
	}

	if count != 2 {
		t.Errorf("deleted count = %d, want 2", count)
	}

	// Verify valid session still exists
	_, err = repo.Get(validSession.ID)
	if err != nil {
		t.Error("valid session was incorrectly deleted")
	}

	// Verify expired sessions are gone
	for _, token := range []string{expiredSession1.ID, expiredSession2.ID} {
		_, err = repo.Get(token)
		if !errors.Is(err, ErrSessionNotFound) {
			t.Errorf("expired session %s still exists", token)
		}
	}
}

func TestSessionRepository_Extend(t *testing.T) {
	store, userID, cleanup := setupSessionTestStore(t)
	defer cleanup()

	repo := NewSessionRepository(store, 24)

	t.Run("extends session successfully", func(t *testing.T) {
		session, _ := repo.Create(userID)
		originalExpiry := session.ExpiresAt

		// Wait a tiny bit to ensure time difference
		time.Sleep(10 * time.Millisecond)

		err := repo.Extend(session.ID)
		if err != nil {
			t.Fatalf("Extend failed: %v", err)
		}

		extended, err := repo.Get(session.ID)
		if err != nil {
			t.Fatalf("Get failed after extend: %v", err)
		}

		if !extended.ExpiresAt.After(originalExpiry) {
			t.Error("expiration was not extended")
		}
	})

	t.Run("returns ErrSessionNotFound for nonexistent session", func(t *testing.T) {
		err := repo.Extend("nonexistent")
		if !errors.Is(err, ErrSessionNotFound) {
			t.Errorf("expected ErrSessionNotFound, got: %v", err)
		}
	})
}

func TestGenerateSessionToken(t *testing.T) {
	tokens := make(map[string]bool)

	// Generate multiple tokens and verify they're unique and correct length
	for i := 0; i < 100; i++ {
		token, err := generateSessionToken()
		if err != nil {
			t.Fatalf("generateSessionToken failed: %v", err)
		}

		if len(token) != 64 {
			t.Errorf("token length = %d, want 64", len(token))
		}

		if tokens[token] {
			t.Error("duplicate token generated")
		}
		tokens[token] = true

		// Verify it's valid hex
		for _, c := range token {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("invalid hex character in token: %c", c)
			}
		}
	}
}
