package film_handler

import "net/http"

func (h *FilmHandler) GetFilms(w http.ResponseWriter, _ *http.Request) {
	films, err := h.filmService.GetFilms()
	if err != nil {
		http.Error(w, "failed to get films", http.StatusInternalServerError)
		return
	}

	_ = films
	w.Write([]byte("Get films"))
}
