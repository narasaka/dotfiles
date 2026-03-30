package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/kubeploy/kubeploy/internal/models"
)

type AuthHandler struct {
	users *models.UserStore
	auth  *AuthMiddleware
}

func NewAuthHandler(users *models.UserStore, auth *AuthMiddleware) *AuthHandler {
	return &AuthHandler{users: users, auth: auth}
}

type setupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AuthHandler) Setup(w http.ResponseWriter, r *http.Request) {
	count, err := h.users.Count()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check users")
		return
	}
	if count > 0 {
		writeError(w, http.StatusBadRequest, "setup already completed")
		return
	}

	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.users.Create(req.Email, string(hash), req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	if err := h.auth.SetSession(w, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.users.GetByEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := h.auth.SetSession(w, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.auth.ClearSession(w)
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) CheckSetup(w http.ResponseWriter, r *http.Request) {
	count, err := h.users.Count()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check users")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"needs_setup": count == 0})
}
