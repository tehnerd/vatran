package middleware

import (
	"context"
	"net/http"

	"github.com/tehnerd/vatran/go/server/models"
)

// contextKey is a type for context keys.
type contextKey string

const (
	// UserIDKey is the context key for the authenticated user ID.
	UserIDKey contextKey = "userID"
)

// AuthResult represents the result of an authentication attempt.
type AuthResult struct {
	// Allowed indicates whether the request is allowed.
	Allowed bool
	// UserID is the authenticated user's identifier.
	UserID string
	// Message is an optional message (e.g., error reason).
	Message string
}

// Authenticator is the interface for authentication implementations.
type Authenticator interface {
	// Authenticate authenticates the request.
	//
	// Parameters:
	//   - r: The HTTP request to authenticate.
	//
	// Returns the authentication result.
	Authenticate(r *http.Request) AuthResult
}

// NoOpAuthenticator is a no-op authenticator that allows all requests.
type NoOpAuthenticator struct{}

// Authenticate implements the Authenticator interface.
// It always allows the request with the user ID set to "anonymous".
//
// Parameters:
//   - r: The HTTP request to authenticate.
//
// Returns an AuthResult that always allows access.
func (n *NoOpAuthenticator) Authenticate(r *http.Request) AuthResult {
	return AuthResult{
		Allowed: true,
		UserID:  "anonymous",
	}
}

// NewNoOpAuthenticator creates a new NoOpAuthenticator.
//
// Returns a new NoOpAuthenticator.
func NewNoOpAuthenticator() *NoOpAuthenticator {
	return &NoOpAuthenticator{}
}

// Auth creates an authentication middleware using the provided authenticator.
//
// Parameters:
//   - auth: The authenticator to use.
//
// Returns a middleware function.
func Auth(auth Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := auth.Authenticate(r)
			if !result.Allowed {
				msg := "unauthorized"
				if result.Message != "" {
					msg = result.Message
				}
				models.WriteError(w, http.StatusUnauthorized, models.NewAPIError(models.CodeUnauthorized, msg))
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, result.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the request context.
//
// Parameters:
//   - r: The HTTP request.
//
// Returns the user ID, or empty string if not found.
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
