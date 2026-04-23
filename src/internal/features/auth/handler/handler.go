package auth_handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	auth_service "github.com/matuha07/kinotower-go/src/internal/features/auth/service"
)

type AuthHandler struct {
	svc auth_service.AuthService
}

func NewAuthHandler(svc auth_service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FIO       string `json:"fio"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Birthday  string `json:"birthday"`
		CountryID *int   `json:"country_id"`
		GenderID  *int   `json:"gender_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.FIO == "" || req.Email == "" || req.Password == "" || req.Birthday == "" || req.GenderID == nil {
		writeError(w, http.StatusBadRequest, "validation error")
		return
	}

	user, err := h.svc.SignUp(domain.SignUpInput{
		FIO:       req.FIO,
		Email:     req.Email,
		Password:  req.Password,
		Birthday:  req.Birthday,
		CountryID: req.CountryID,
		GenderID:  req.GenderID,
	})
	if err != nil {
		if errors.Is(err, auth_service.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation error")
		return
	}

	tokens, err := h.svc.SignIn(domain.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, auth_service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, domain.InvalidAuthResponse{Status: "invalid", Message: "Wrong email or password"})
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}
