package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Session represents an authenticated user session.
type Session struct {
	// ID is the session token (256-bit random hex string).
	ID string
	// UserID is the ID of the authenticated user.
	UserID int64
	// CreatedAt is when the session was created.
	CreatedAt time.Time
	// ExpiresAt is when the session expires.
	ExpiresAt time.Time
}

// ErrSessionNotFound is returned when a session is not found or expired.
var ErrSessionNotFound = errors.New("session not found")

// ErrSessionExpired is returned when a session has expired.
var ErrSessionExpired = errors.New("session expired")

// SessionRepository handles session persistence operations.
type SessionRepository struct {
	store          *Store
	sessionTimeout time.Duration
}

// NewSessionRepository creates a new SessionRepository.
//
// Parameters:
//   - store: The database store.
//   - sessionTimeoutHours: The session timeout in hours.
//
// Returns a new SessionRepository instance.
func NewSessionRepository(store *Store, sessionTimeoutHours int) *SessionRepository {
	if sessionTimeoutHours <= 0 {
		sessionTimeoutHours = 24
	}
	return &SessionRepository{
		store:          store,
		sessionTimeout: time.Duration(sessionTimeoutHours) * time.Hour,
	}
}

// Create creates a new session for the given user.
//
// Parameters:
//   - userID: The ID of the user to create a session for.
//
// Returns the created Session or an error.
func (r *SessionRepository) Create(userID int64) (*Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(r.sessionTimeout)

	_, err = r.store.DB().Exec(
		"INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)",
		token, userID, now, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

// Get retrieves a session by its token.
// Returns ErrSessionNotFound if the session doesn't exist.
// Returns ErrSessionExpired if the session has expired (and deletes it).
//
// Parameters:
//   - token: The session token to look up.
//
// Returns the Session or an error.
func (r *SessionRepository) Get(token string) (*Session, error) {
	var session Session
	err := r.store.DB().QueryRow(
		"SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?",
		token,
	).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		// Delete expired session
		r.Delete(token)
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// Delete removes a session by its token.
//
// Parameters:
//   - token: The session token to delete.
//
// Returns an error if deletion fails.
func (r *SessionRepository) Delete(token string) error {
	_, err := r.store.DB().Exec("DELETE FROM sessions WHERE id = ?", token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteByUserID removes all sessions for a given user.
//
// Parameters:
//   - userID: The user ID whose sessions should be deleted.
//
// Returns an error if deletion fails.
func (r *SessionRepository) DeleteByUserID(userID int64) error {
	_, err := r.store.DB().Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

// DeleteExpired removes all expired sessions from the database.
//
// Returns the number of deleted sessions or an error.
func (r *SessionRepository) DeleteExpired() (int64, error) {
	result, err := r.store.DB().Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return count, nil
}

// Extend extends a session's expiration time.
//
// Parameters:
//   - token: The session token to extend.
//
// Returns an error if the update fails.
func (r *SessionRepository) Extend(token string) error {
	expiresAt := time.Now().Add(r.sessionTimeout)
	result, err := r.store.DB().Exec(
		"UPDATE sessions SET expires_at = ? WHERE id = ?",
		expiresAt, token,
	)
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return ErrSessionNotFound
	}

	return nil
}

// generateSessionToken generates a cryptographically secure 256-bit session token.
func generateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
