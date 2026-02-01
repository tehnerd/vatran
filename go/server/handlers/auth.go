package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/tehnerd/vatran/go/server/auth"
	"github.com/tehnerd/vatran/go/server/models"
)

// loginPageTemplate is the inline HTML template for the login page.
const loginPageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Katran</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .login-container {
            background: #fff;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
            width: 100%;
            max-width: 400px;
        }
        h1 {
            text-align: center;
            color: #1a1a2e;
            margin-bottom: 8px;
            font-size: 28px;
        }
        .subtitle {
            text-align: center;
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            color: #333;
            font-weight: 500;
        }
        input[type="text"],
        input[type="password"] {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #e1e1e1;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        input[type="text"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: #4a90d9;
        }
        button {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #4a90d9 0%, #357abd 100%);
            color: #fff;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.1s, box-shadow 0.2s;
        }
        button:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(74, 144, 217, 0.4);
        }
        button:active {
            transform: translateY(0);
        }
        .error {
            background: #fee;
            color: #c00;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 20px;
            text-align: center;
            font-size: 14px;
        }
        .logo {
            text-align: center;
            margin-bottom: 20px;
            font-size: 48px;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="logo">&#9881;</div>
        <h1>Katran</h1>
        <p class="subtitle">L4 Load Balancer</p>
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
        <form method="POST" action="/login{{if .Redirect}}?redirect={{.Redirect}}{{end}}">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autofocus>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">Sign In</button>
        </form>
    </div>
</body>
</html>`

// LoginRequest represents a JSON login request.
type LoginRequest struct {
	// Username is the username to authenticate.
	Username string `json:"username"`
	// Password is the password to authenticate.
	Password string `json:"password"`
}

// LoginResponse represents a successful login response.
type LoginResponse struct {
	// Success indicates the login was successful.
	Success bool `json:"success"`
	// Username is the authenticated username.
	Username string `json:"username"`
}

// AuthHandler handles authentication-related HTTP endpoints.
type AuthHandler struct {
	authenticator *auth.BasicAuthenticator
	loginTemplate *template.Template
}

// NewAuthHandler creates a new AuthHandler.
//
// Parameters:
//   - authenticator: The BasicAuthenticator instance.
//
// Returns a new AuthHandler instance.
func NewAuthHandler(authenticator *auth.BasicAuthenticator) *AuthHandler {
	tmpl := template.Must(template.New("login").Parse(loginPageTemplate))
	return &AuthHandler{
		authenticator: authenticator,
		loginTemplate: tmpl,
	}
}

// HandleLogin handles GET and POST /login.
// GET renders the login page.
// POST validates credentials and creates a session.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.renderLoginPage(w, r, "")
	case http.MethodPost:
		h.handleLoginSubmit(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// renderLoginPage renders the login page with an optional error message.
func (h *AuthHandler) renderLoginPage(w http.ResponseWriter, r *http.Request, errorMsg string) {
	redirect := r.URL.Query().Get("redirect")
	data := struct {
		Error    string
		Redirect string
	}{
		Error:    errorMsg,
		Redirect: redirect,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.loginTemplate.Execute(w, data)
}

// handleLoginSubmit processes login form submissions.
func (h *AuthHandler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	var username, password string
	redirect := r.URL.Query().Get("redirect")

	// Check content type to determine if this is JSON or form submission
	contentType := r.Header.Get("Content-Type")
	isJSON := strings.Contains(contentType, "application/json")

	if isJSON {
		// Parse JSON request
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			models.WriteError(w, http.StatusBadRequest,
				models.NewInvalidRequestError("invalid request body: "+err.Error()))
			return
		}
		username = req.Username
		password = req.Password
	} else {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			h.renderLoginPage(w, r, "Invalid form data")
			return
		}
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

	// Validate credentials
	if username == "" || password == "" {
		if isJSON {
			models.WriteError(w, http.StatusBadRequest,
				models.NewInvalidRequestError("username and password are required"))
		} else {
			h.renderLoginPage(w, r, "Username and password are required")
		}
		return
	}

	// Authenticate
	token, expiresAt, err := h.authenticator.Login(username, password)
	if err != nil {
		if isJSON {
			models.WriteError(w, http.StatusUnauthorized,
				models.NewAPIError(models.CodeUnauthorized, "Invalid username or password"))
		} else {
			h.renderLoginPage(w, r, "Invalid username or password")
		}
		return
	}

	// Set session cookie
	h.authenticator.SetSessionCookie(w, token, expiresAt)

	if isJSON {
		models.WriteSuccess(w, LoginResponse{
			Success:  true,
			Username: username,
		})
	} else {
		// Redirect to original URL or home
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}

// HandleLogout handles POST /logout.
// Invalidates the session and clears the cookie.
//
// Parameters:
//   - w: The http.ResponseWriter to write to.
//   - r: The HTTP request.
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session cookie
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err == nil && cookie.Value != "" {
		// Invalidate session
		h.authenticator.Logout(cookie.Value)
	}

	// Clear cookie
	h.authenticator.ClearSessionCookie(w)

	// Check if JSON response is expected
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		models.WriteSuccess(w, map[string]bool{"logged_out": true})
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
