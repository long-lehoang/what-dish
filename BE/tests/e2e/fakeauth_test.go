package e2e_test

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// fakeAuthServer mimics the Supabase Auth REST API for E2E tests.
// It listens on a fixed port so the app container can reach it via host.docker.internal.
type fakeAuthServer struct {
	server *http.Server
	ln     net.Listener
	mu     sync.Mutex
	users  map[string]fakeUser // email -> user
}

type fakeUser struct {
	ID       uuid.UUID
	Email    string
	Name     string
	Password string
}

// startFakeAuthServerOnAddr starts a fake Supabase Auth HTTP server on the given address (e.g. ":9999").
func startFakeAuthServerOnAddr(addr string) *fakeAuthServer {
	fa := &fakeAuthServer{
		users: make(map[string]fakeUser),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/v1/signup", fa.handleSignup)
	mux.HandleFunc("/auth/v1/token", fa.handleToken)
	mux.HandleFunc("/auth/v1/user", fa.handleGetUser)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic("fakeauth: listen " + addr + ": " + err.Error())
	}
	fa.ln = ln
	fa.server = &http.Server{Handler: mux}

	go fa.server.Serve(ln)
	return fa
}

func (fa *fakeAuthServer) Close() {
	if fa.server != nil {
		fa.server.Close()
	}
}

func (fa *fakeAuthServer) handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Email    string            `json:"email"`
		Password string            `json:"password"`
		Data     map[string]string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthError(w, http.StatusBadRequest, "invalid_request", "invalid body")
		return
	}

	fa.mu.Lock()
	if _, exists := fa.users[body.Email]; exists {
		fa.mu.Unlock()
		writeAuthError(w, http.StatusConflict, "user_already_exists", "user already registered")
		return
	}

	u := fakeUser{
		ID:       uuid.New(),
		Email:    body.Email,
		Name:     body.Data["name"],
		Password: body.Password,
	}
	fa.users[body.Email] = u
	fa.mu.Unlock()

	writeAuthResponse(w, http.StatusOK, u)
}

func (fa *fakeAuthServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	grantType := r.URL.Query().Get("grant_type")

	switch grantType {
	case "password":
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeAuthError(w, http.StatusBadRequest, "invalid_request", "invalid body")
			return
		}

		fa.mu.Lock()
		u, exists := fa.users[body.Email]
		fa.mu.Unlock()

		if !exists || u.Password != body.Password {
			writeAuthError(w, http.StatusUnauthorized, "invalid_grant", "invalid login credentials")
			return
		}

		writeAuthResponse(w, http.StatusOK, u)

	case "refresh_token":
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeAuthError(w, http.StatusBadRequest, "invalid_request", "invalid body")
			return
		}

		userID := uuid.New()
		if strings.HasPrefix(body.RefreshToken, "refresh-") {
			if parsed, err := uuid.Parse(strings.TrimPrefix(body.RefreshToken, "refresh-")); err == nil {
				userID = parsed
			}
		}

		u := fakeUser{ID: userID, Email: "refreshed@test.com", Name: "Refreshed User"}
		writeAuthResponse(w, http.StatusOK, u)

	default:
		writeAuthError(w, http.StatusBadRequest, "unsupported_grant_type", "unsupported grant type")
	}
}

func (fa *fakeAuthServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")

	if !strings.HasPrefix(token, "test-token-") {
		http.Error(w, `{"error":"invalid_token"}`, http.StatusUnauthorized)
		return
	}

	userIDStr := strings.TrimPrefix(token, "test-token-")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid_token"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":    userID.String(),
		"email": "testuser@example.com",
		"user_metadata": map[string]string{
			"name": "Test User",
		},
	})
}

func writeAuthResponse(w http.ResponseWriter, status int, u fakeUser) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"access_token":  "test-token-" + u.ID.String(),
		"refresh_token": "refresh-" + u.ID.String(),
		"expires_in":    3600,
		"user": map[string]any{
			"id":    u.ID.String(),
			"email": u.Email,
			"user_metadata": map[string]string{
				"name": u.Name,
			},
		},
	})
}

func writeAuthError(w http.ResponseWriter, status int, errCode, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             errCode,
		"error_description": desc,
	})
}
