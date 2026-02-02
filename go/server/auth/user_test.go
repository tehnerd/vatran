package auth

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func setupUserTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "user_test")
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

func TestNewUserRepository(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	t.Run("valid bcrypt cost", func(t *testing.T) {
		repo := NewUserRepository(store, 10)
		if repo == nil {
			t.Error("NewUserRepository returned nil")
		}
		if repo.bcryptCost != 10 {
			t.Errorf("bcryptCost = %d, want 10", repo.bcryptCost)
		}
	})

	t.Run("invalid bcrypt cost uses default", func(t *testing.T) {
		repo := NewUserRepository(store, 0)
		if repo.bcryptCost != bcrypt.DefaultCost {
			t.Errorf("bcryptCost = %d, want %d", repo.bcryptCost, bcrypt.DefaultCost)
		}

		repo = NewUserRepository(store, 50)
		if repo.bcryptCost != bcrypt.DefaultCost {
			t.Errorf("bcryptCost = %d, want %d", repo.bcryptCost, bcrypt.DefaultCost)
		}
	})
}

func TestUserRepository_Create(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost) // Use MinCost for faster tests

	t.Run("creates user successfully", func(t *testing.T) {
		user, err := repo.Create("testuser", "password123")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if user.ID == 0 {
			t.Error("user ID is 0")
		}
		if user.Username != "testuser" {
			t.Errorf("Username = %q, want %q", user.Username, "testuser")
		}
		if user.PasswordHash == "" {
			t.Error("PasswordHash is empty")
		}
		if user.PasswordHash == "password123" {
			t.Error("PasswordHash contains plaintext password")
		}

		// Verify password hash is valid bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
		if err != nil {
			t.Error("password hash does not match original password")
		}
	})

	t.Run("returns ErrUserExists for duplicate username", func(t *testing.T) {
		_, err := repo.Create("duplicate", "password1")
		if err != nil {
			t.Fatalf("first Create failed: %v", err)
		}

		_, err = repo.Create("duplicate", "password2")
		if !errors.Is(err, ErrUserExists) {
			t.Errorf("expected ErrUserExists, got: %v", err)
		}
	})
}

func TestUserRepository_Authenticate(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost)

	// Create test user
	_, err := repo.Create("authuser", "correctpassword")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("valid credentials", func(t *testing.T) {
		user, err := repo.Authenticate("authuser", "correctpassword")
		if err != nil {
			t.Fatalf("Authenticate failed: %v", err)
		}
		if user.Username != "authuser" {
			t.Errorf("Username = %q, want %q", user.Username, "authuser")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := repo.Authenticate("authuser", "wrongpassword")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("nonexistent user", func(t *testing.T) {
		_, err := repo.Authenticate("noexist", "password")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost)

	// Create test user
	created, err := repo.Create("findme", "password")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("finds existing user", func(t *testing.T) {
		user, err := repo.FindByUsername("findme")
		if err != nil {
			t.Fatalf("FindByUsername failed: %v", err)
		}
		if user.ID != created.ID {
			t.Errorf("ID = %d, want %d", user.ID, created.ID)
		}
		if user.Username != "findme" {
			t.Errorf("Username = %q, want %q", user.Username, "findme")
		}
	})

	t.Run("returns ErrUserNotFound for nonexistent user", func(t *testing.T) {
		_, err := repo.FindByUsername("noexist")
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got: %v", err)
		}
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost)

	// Create test user
	created, err := repo.Create("findbyid", "password")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("finds existing user", func(t *testing.T) {
		user, err := repo.FindByID(created.ID)
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}
		if user.Username != "findbyid" {
			t.Errorf("Username = %q, want %q", user.Username, "findbyid")
		}
	})

	t.Run("returns ErrUserNotFound for nonexistent ID", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got: %v", err)
		}
	})
}

func TestUserRepository_UserCount(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost)

	t.Run("returns 0 for empty database", func(t *testing.T) {
		count, err := repo.UserCount()
		if err != nil {
			t.Fatalf("UserCount failed: %v", err)
		}
		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("returns correct count", func(t *testing.T) {
		_, _ = repo.Create("user1", "pass")
		_, _ = repo.Create("user2", "pass")
		_, _ = repo.Create("user3", "pass")

		count, err := repo.UserCount()
		if err != nil {
			t.Fatalf("UserCount failed: %v", err)
		}
		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	store, cleanup := setupUserTestStore(t)
	defer cleanup()

	repo := NewUserRepository(store, bcrypt.MinCost)

	// Create test user
	user, err := repo.Create("updatepw", "oldpassword")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("updates password successfully", func(t *testing.T) {
		err := repo.UpdatePassword(user.ID, "newpassword")
		if err != nil {
			t.Fatalf("UpdatePassword failed: %v", err)
		}

		// Verify old password no longer works
		_, err = repo.Authenticate("updatepw", "oldpassword")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Error("old password still works")
		}

		// Verify new password works
		_, err = repo.Authenticate("updatepw", "newpassword")
		if err != nil {
			t.Errorf("new password authentication failed: %v", err)
		}
	})

	t.Run("returns ErrUserNotFound for nonexistent user", func(t *testing.T) {
		err := repo.UpdatePassword(99999, "newpass")
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got: %v", err)
		}
	})
}
