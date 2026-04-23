package film_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/matuha07/kinotower-go/src/internal/core/domain"
	film_service "github.com/matuha07/kinotower-go/src/internal/features/films/service"
)

type FilmHandler struct {
	filmService film_service.FilmService
}

func NewFilmHandler(filmService film_service.FilmService) *FilmHandler {
	return &FilmHandler{filmService: filmService}
}

func (h *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	films, err := h.filmService.GetFilms(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get films")
		return
	}
	writeJSON(w, http.StatusOK, films)
}

func (h *FilmHandler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	film, err := h.filmService.GetFilmByID(id)
	if err != nil {
		if errors.Is(err, film_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Film not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get film")
		return
	}

	writeJSON(w, http.StatusOK, film)
}

func (h *FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	var req createFilmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Name == "" || req.Duration <= 0 || req.YearOfIssue <= 0 {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	film, err := h.filmService.CreateFilm(req.toDomain())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create film")
		return
	}

	writeJSON(w, http.StatusCreated, film)
}

func (h *FilmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateFilmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	film, err := h.filmService.UpdateFilm(id, req.toDomain())
	if err != nil {
		if errors.Is(err, film_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "film not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update film")
		return
	}

	writeJSON(w, http.StatusOK, film)
}

func (h *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.filmService.DeleteFilm(id); err != nil {
		if errors.Is(err, film_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "film not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete film")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FilmHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	categories, err := h.filmService.GetCategories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get categories")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": categories})
}

func (h *FilmHandler) GetCountries(w http.ResponseWriter, _ *http.Request) {
	countries, err := h.filmService.GetCountries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get countries")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"countries": countries})
}

func (h *FilmHandler) GetGenders(w http.ResponseWriter, _ *http.Request) {
	genders, err := h.filmService.GetGenders()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get genders")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"genders": genders})
}

func (h *FilmHandler) GetFilmReviews(w http.ResponseWriter, r *http.Request) {
	filmID, err := parseID(chi.URLParam(r, "film-id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid film-id")
		return
	}

	reviews, err := h.filmService.GetFilmReviews(filmID)
	if err != nil {
		if errors.Is(err, film_service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Film not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get reviews")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reviews": reviews})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseID(raw string) (int, error) {
	return strconv.Atoi(raw)
}

func parseFilter(r *http.Request) (domain.Filter, error) {
	q := r.URL.Query()
	filter := domain.Filter{Page: 1, Size: 10, SortBy: "name", SortDir: "asc"}

	if raw := q.Get("page"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 {
			return filter, errors.New("invalid page")
		}
		filter.Page = v
	}
	if raw := q.Get("size"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 {
			return filter, errors.New("invalid size")
		}
		filter.Size = v
	}
	if raw := q.Get("sortBy"); raw != "" {
		if raw != "name" && raw != "year" && raw != "rating" {
			return filter, errors.New("invalid sortBy")
		}
		filter.SortBy = raw
	}
	if raw := q.Get("sortDir"); raw != "" {
		if raw != "asc" && raw != "desc" {
			return filter, errors.New("invalid sortDir")
		}
		filter.SortDir = raw
	}
	if raw := q.Get("category"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 {
			return filter, errors.New("invalid category")
		}
		filter.CategoryID = v
	}
	if raw := q.Get("country"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 {
			return filter, errors.New("invalid country")
		}
		filter.CountryID = v
	}
	filter.Search = q.Get("search")

	return filter, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}
