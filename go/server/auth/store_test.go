package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	// Use a temp directory for test databases
	tmpDir, err := os.MkdirTemp("", "auth_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("creates new database", func(t *testing.T) {
		dbPath := filepath.Join(tmpDir, "test1.db")

		store, err := NewStore(dbPath)
		if err != nil {
			t.Fatalf("NewStore failed: %v", err)
		}
		defer store.Close()

		// Verify database file was created
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Error("database file was not created")
		}

		// Verify DB is accessible
		if store.DB() == nil {
			t.Error("DB() returned nil")
		}

		// Verify path is correct
		if store.Path() != dbPath {
			t.Errorf("Path() = %q, want %q", store.Path(), dbPath)
		}
	})

	t.Run("creates nested directory", func(t *testing.T) {
		dbPath := filepath.Join(tmpDir, "nested", "dir", "test2.db")

		store, err := NewStore(dbPath)
		if err != nil {
			t.Fatalf("NewStore failed: %v", err)
		}
		defer store.Close()

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Error("database file was not created in nested directory")
		}
	})

	t.Run("schema is applied", func(t *testing.T) {
		dbPath := filepath.Join(tmpDir, "test3.db")

		store, err := NewStore(dbPath)
		if err != nil {
			t.Fatalf("NewStore failed: %v", err)
		}
		defer store.Close()

		// Verify users table exists
		var name string
		err = store.DB().QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name='users'",
		).Scan(&name)
		if err != nil {
			t.Errorf("users table not found: %v", err)
		}

		// Verify sessions table exists
		err = store.DB().QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name='sessions'",
		).Scan(&name)
		if err != nil {
			t.Errorf("sessions table not found: %v", err)
		}
	})

	t.Run("reopens existing database", func(t *testing.T) {
		dbPath := filepath.Join(tmpDir, "test4.db")

		// Create initial store
		store1, err := NewStore(dbPath)
		if err != nil {
			t.Fatalf("NewStore failed: %v", err)
		}

		// Insert test data
		_, err = store1.DB().Exec(
			"INSERT INTO users (username, password_hash) VALUES (?, ?)",
			"testuser", "testhash",
		)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
		store1.Close()

		// Reopen store
		store2, err := NewStore(dbPath)
		if err != nil {
			t.Fatalf("NewStore failed on reopen: %v", err)
		}
		defer store2.Close()

		// Verify data persisted
		var username string
		err = store2.DB().QueryRow(
			"SELECT username FROM users WHERE username = ?", "testuser",
		).Scan(&username)
		if err != nil {
			t.Errorf("data not persisted: %v", err)
		}
	})
}

func TestStore_Close(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "auth_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "close_test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// Close should not error
	if err := store.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Double close should not error
	if err := store.Close(); err != nil {
		t.Errorf("second Close() returned error: %v", err)
	}
}

func TestStore_ForeignKeysEnabled(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "auth_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "fk_test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Verify foreign keys are enabled
	var fkEnabled int
	err = store.DB().QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("failed to query foreign_keys pragma: %v", err)
	}
	if fkEnabled != 1 {
		t.Errorf("foreign_keys = %d, want 1", fkEnabled)
	}
}
