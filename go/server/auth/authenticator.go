package auth

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/tehnerd/vatran/go/server/middleware"
)

const (
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "katran_session"
)

// BasicAuthenticator implements middleware.Authenticator using session-based auth.
type BasicAuthenticator struct {
	userRepo       *UserRepository
	sessionRepo    *SessionRepository
	store          *Store
	allowLocalhost bool
	exemptPaths    []string
	isTLS          bool
}

// NewBasicAuthenticator creates a new BasicAuthenticator.
//
// Parameters:
//   - store: The database store.
//   - bcryptCost: The bcrypt cost factor for password hashing.
//   - sessionTimeoutHours: The session timeout in hours.
//   - allowLocalhost: Whether to bypass auth for localhost requests.
//   - exemptPaths: List of paths to exempt from authentication (supports * suffix for prefix matching).
//   - isTLS: Whether the server is using TLS (for secure cookies).
//
// Returns a new BasicAuthenticator instance.
func NewBasicAuthenticator(
	store *Store,
	bcryptCost int,
	sessionTimeoutHours int,
	allowLocalhost bool,
	exemptPaths []string,
	isTLS bool,
) *BasicAuthenticator {
	// Always exempt /login and /health
	defaultExempts := []string{"/login", "/health"}
	allExempts := append(defaultExempts, exemptPaths...)

	return &BasicAuthenticator{
		userRepo:       NewUserRepository(store, bcryptCost),
		sessionRepo:    NewSessionRepository(store, sessionTimeoutHours),
		store:          store,
		allowLocalhost: allowLocalhost,
		exemptPaths:    allExempts,
		isTLS:          isTLS,
	}
}

// Authenticate implements the middleware.Authenticator interface.
//
// Parameters:
//   - r: The HTTP request to authenticate.
//
// Returns the authentication result.
func (a *BasicAuthenticator) Authenticate(r *http.Request) middleware.AuthResult {
	// Check if path is exempt
	if a.isExemptPath(r.URL.Path) {
		return middleware.AuthResult{
			Allowed: true,
			UserID:  "anonymous",
		}
	}

	// Check localhost bypass
	if a.allowLocalhost && a.isLocalhost(r) {
		return middleware.AuthResult{
			Allowed: true,
			UserID:  "localhost",
		}
	}

	// Check session cookie
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return middleware.AuthResult{
			Allowed: false,
			Message: "authentication required",
		}
	}

	session, err := a.sessionRepo.Get(cookie.Value)
	if err != nil {
		return middleware.AuthResult{
			Allowed: false,
			Message: "invalid or expired session",
		}
	}

	// Find user for the session
	user, err := a.userRepo.FindByID(session.UserID)
	if err != nil {
		return middleware.AuthResult{
			Allowed: false,
			Message: "user not found",
		}
	}

	return middleware.AuthResult{
		Allowed: true,
		UserID:  user.Username,
	}
}

// Login authenticates a user and creates a session.
//
// Parameters:
//   - username: The username.
//   - password: The plaintext password.
//
// Returns the session token and expiration time, or an error.
func (a *BasicAuthenticator) Login(username, password string) (string, time.Time, error) {
	user, err := a.userRepo.Authenticate(username, password)
	if err != nil {
		return "", time.Time{}, err
	}

	session, err := a.sessionRepo.Create(user.ID)
	if err != nil {
		return "", time.Time{}, err
	}

	return session.ID, session.ExpiresAt, nil
}

// Logout invalidates a session.
//
// Parameters:
//   - token: The session token to invalidate.
//
// Returns an error if logout fails.
func (a *BasicAuthenticator) Logout(token string) error {
	return a.sessionRepo.Delete(token)
}

// CreateUser creates a new user.
//
// Parameters:
//   - username: The username for the new user.
//   - password: The plaintext password.
//
// Returns the created User or an error.
func (a *BasicAuthenticator) CreateUser(username, password string) (*User, error) {
	return a.userRepo.Create(username, password)
}

// UserCount returns the total number of users.
//
// Returns the count or an error.
func (a *BasicAuthenticator) UserCount() (int64, error) {
	return a.userRepo.UserCount()
}

// Close closes the underlying store.
//
// Returns an error if closing fails.
func (a *BasicAuthenticator) Close() error {
	return a.store.Close()
}

// IsTLS returns whether TLS is enabled.
//
// Returns true if TLS is enabled.
func (a *BasicAuthenticator) IsTLS() bool {
	return a.isTLS
}

// CleanupExpiredSessions removes expired sessions from the database.
//
// Returns the number of deleted sessions or an error.
func (a *BasicAuthenticator) CleanupExpiredSessions() (int64, error) {
	return a.sessionRepo.DeleteExpired()
}

// isExemptPath checks if a path is exempt from authentication.
func (a *BasicAuthenticator) isExemptPath(path string) bool {
	for _, exempt := range a.exemptPaths {
		if exempt == path {
			return true
		}
		// Support wildcard prefix matching (e.g., "/static/*")
		if strings.HasSuffix(exempt, "/*") {
			prefix := strings.TrimSuffix(exempt, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}

// isLocalhost checks if the request is from localhost.
func (a *BasicAuthenticator) isLocalhost(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	return ip.IsLoopback()
}

// SetSessionCookie sets the session cookie on the response.
//
// Parameters:
//   - w: The http.ResponseWriter.
//   - token: The session token.
//   - expiresAt: When the cookie should expire.
func (a *BasicAuthenticator) SetSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   a.isTLS,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearSessionCookie clears the session cookie.
//
// Parameters:
//   - w: The http.ResponseWriter.
func (a *BasicAuthenticator) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.isTLS,
		SameSite: http.SameSiteStrictMode,
	})
}
