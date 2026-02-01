package middleware

import (
	"context"
	"net/http"
	"net/url"
	"strings"

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

				// Check if this is an API request or a browser request
				if isAPIRequest(r) {
					// API requests get JSON 401 response
					models.WriteError(w, http.StatusUnauthorized, models.NewAPIError(models.CodeUnauthorized, msg))
				} else {
					// Browser requests get redirected to login
					redirectToLogin(w, r)
				}
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, result.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isAPIRequest determines if the request is an API request vs a browser request.
// API requests include:
// - Requests to /api/* paths
// - Requests with Accept: application/json
// - Requests with X-Requested-With: XMLHttpRequest (AJAX)
//
// Parameters:
//   - r: The HTTP request.
//
// Returns true if this is an API request.
func isAPIRequest(r *http.Request) bool {
	// Check path prefix
	if strings.HasPrefix(r.URL.Path, "/api/") {
		return true
	}

	// Check Accept header for JSON
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}

	// Check X-Requested-With header (AJAX requests)
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}

	return false
}

// redirectToLogin redirects the request to the login page with a redirect parameter.
//
// Parameters:
//   - w: The http.ResponseWriter.
//   - r: The HTTP request.
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	// Build redirect URL with original request path
	loginURL := "/login"
	if r.URL.Path != "/" && r.URL.Path != "/login" {
		redirectParam := r.URL.Path
		if r.URL.RawQuery != "" {
			redirectParam += "?" + r.URL.RawQuery
		}
		loginURL = "/login?redirect=" + url.QueryEscape(redirectParam)
	}
	http.Redirect(w, r, loginURL, http.StatusFound)
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
