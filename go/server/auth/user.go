package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	// ID is the unique identifier for the user.
	ID int64
	// Username is the user's login name.
	Username string
	// PasswordHash is the bcrypt hash of the user's password.
	PasswordHash string
	// CreatedAt is when the user was created.
	CreatedAt time.Time
	// UpdatedAt is when the user was last updated.
	UpdatedAt time.Time
}

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrUserExists is returned when attempting to create a user that already exists.
var ErrUserExists = errors.New("user already exists")

// ErrInvalidCredentials is returned when authentication fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// UserRepository handles user persistence operations.
type UserRepository struct {
	store      *Store
	bcryptCost int
}

// NewUserRepository creates a new UserRepository.
//
// Parameters:
//   - store: The database store.
//   - bcryptCost: The bcrypt cost factor (recommended: 12).
//
// Returns a new UserRepository instance.
func NewUserRepository(store *Store, bcryptCost int) *UserRepository {
	if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
		bcryptCost = bcrypt.DefaultCost
	}
	return &UserRepository{
		store:      store,
		bcryptCost: bcryptCost,
	}
}

// Create creates a new user with the given username and password.
//
// Parameters:
//   - username: The username for the new user.
//   - password: The plaintext password (will be hashed).
//
// Returns the created User or an error.
func (r *UserRepository) Create(username, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), r.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := r.store.DB().Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username, string(hash),
	)
	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}

	return r.FindByID(id)
}

// Authenticate verifies the username and password combination.
//
// Parameters:
//   - username: The username to authenticate.
//   - password: The plaintext password to verify.
//
// Returns the authenticated User or an error.
func (r *UserRepository) Authenticate(username, password string) (*User, error) {
	user, err := r.FindByUsername(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// FindByUsername finds a user by their username.
//
// Parameters:
//   - username: The username to search for.
//
// Returns the User or ErrUserNotFound.
func (r *UserRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.store.DB().QueryRow(
		"SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// FindByID finds a user by their ID.
//
// Parameters:
//   - id: The user ID to search for.
//
// Returns the User or ErrUserNotFound.
func (r *UserRepository) FindByID(id int64) (*User, error) {
	var user User
	err := r.store.DB().QueryRow(
		"SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// UserCount returns the total number of users in the database.
//
// Returns the count or an error.
func (r *UserRepository) UserCount() (int64, error) {
	var count int64
	err := r.store.DB().QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// UpdatePassword updates a user's password.
//
// Parameters:
//   - id: The user ID.
//   - password: The new plaintext password.
//
// Returns an error if the update fails.
func (r *UserRepository) UpdatePassword(id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), r.bcryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := r.store.DB().Exec(
		"UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		string(hash), id,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// isUniqueConstraintError checks if the error is a SQLite unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// SQLite unique constraint error contains "UNIQUE constraint failed"
	return err.Error() == "UNIQUE constraint failed: users.username" ||
		(len(err.Error()) > 0 && err.Error()[0:4] == "UNIQ")
}
