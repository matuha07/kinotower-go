package user_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	core_middleware "github.com/matuha07/kinotower-go/src/internal/core/middleware"
	user_service "github.com/matuha07/kinotower-go/src/internal/features/users/service"
)

type UserHandler struct {
	svc user_service.UserService
}

func NewUserHandler(svc user_service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	user, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := core_middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		FIO      *string `json:"fio"`
		Email    *string `json:"email"`
		Birthday *string `json:"birthday"`
		GenderID *int    `json:"gender_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	_, err := h.svc.Update(userID, domain.UserUpdate{
		FIO:      req.FIO,
		Email:    req.Email,
		Birthday: req.Birthday,
		GenderID: req.GenderID,
	})
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, user_service.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := core_middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.svc.Delete(userID); err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}

	var req struct {
		FilmID  int    `json:"film_id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.FilmID <= 0 || len(req.Message) < 4 || len(req.Message) > 1024 {
		writeError(w, http.StatusBadRequest, "validation error")
		return
	}

	review, err := h.svc.CreateReview(userID, domain.ReviewCreate{FilmID: req.FilmID, Message: req.Message})
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, user_service.ErrFilmNotFound) {
			writeError(w, http.StatusNotFound, "Film not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, review)
}

func (h *UserHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}
	reviews, err := h.svc.GetReviews(userID)
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reviews": reviews})
}

func (h *UserHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}
	reviewID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.svc.DeleteReview(userID, reviewID)
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, user_service.ErrReviewNotFound) {
			writeError(w, http.StatusNotFound, "Review not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}

	var req struct {
		FilmID int `json:"film_id"`
		Ball   int `json:"ball"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.FilmID <= 0 || req.Ball < 1 || req.Ball > 5 {
		writeError(w, http.StatusBadRequest, "validation error")
		return
	}

	rating, err := h.svc.CreateRating(userID, domain.RatingCreate{FilmID: req.FilmID, Ball: req.Ball})
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, user_service.ErrFilmNotFound) {
			writeError(w, http.StatusNotFound, "Film not found")
			return
		}
		if errors.Is(err, user_service.ErrScoreExists) {
			writeJSON(w, http.StatusUnauthorized, domain.InvalidAuthResponse{Status: "invalid", Message: "Score exist"})
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, rating)
}

func (h *UserHandler) GetRatings(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}
	ratings, err := h.svc.GetRatings(userID)
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ratings": ratings})
}

func (h *UserHandler) DeleteRating(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user-id")
		return
	}
	ratingID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.svc.DeleteRating(userID, ratingID)
	if err != nil {
		if errors.Is(err, user_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, user_service.ErrRatingNotFound) {
			writeError(w, http.StatusNotFound, "Rating not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}

func parseUserIDParam(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "user-id"))
}
