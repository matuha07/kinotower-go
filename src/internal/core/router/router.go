package core_router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	core_logger "github.com/matuha07/kinotower-go/src/internal/core/logger"
	core_middleware "github.com/matuha07/kinotower-go/src/internal/core/middleware"
	auth_handler "github.com/matuha07/kinotower-go/src/internal/features/auth/handler"
	film_handler "github.com/matuha07/kinotower-go/src/internal/features/films/handler"
	user_handler "github.com/matuha07/kinotower-go/src/internal/features/users/handler"
)

type Router struct {
	mux         *chi.Mux
	filmHandler *film_handler.FilmHandler
	userHandler *user_handler.UserHandler
	authHandler *auth_handler.AuthHandler
	jwtSecret   string
}

func NewRouter(
	filmHandler *film_handler.FilmHandler,
	userHandler *user_handler.UserHandler,
	authHandler *auth_handler.AuthHandler,
	jwtSecret string,
) *Router {
	r := &Router{
		mux:         chi.NewRouter(),
		filmHandler: filmHandler,
		userHandler: userHandler,
		authHandler: authHandler,
		jwtSecret:   jwtSecret,
	}

	r.mux.Use(middleware.RequestID)
	r.mux.Use(core_logger.HTTPMiddleware(slog.Default()))
	r.mux.Use(middleware.Recoverer)

	r.mux.Route("/api/v1", func(rl chi.Router) {
		rl.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("Welcome to the Kinotower API!"))
		})

		// Public
		rl.Mount("/films", r.filmRoutes())
		rl.Get("/categories", r.filmHandler.GetCategories)
		rl.Get("/countries", r.filmHandler.GetCountries)
		rl.Mount("/genders", r.genderRoutes())
		rl.Post("/auth/signup", r.authHandler.SignUp)
		rl.Post("/auth/signin", r.authHandler.SignIn)

		// Protected
		rl.Group(func(rl chi.Router) {
			rl.Use(core_middleware.JWT(jwtSecret))
			rl.Post("/auth/signout", r.authHandler.SignOut)
			rl.Mount("/users", r.userRoutes())
		})
	})

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
