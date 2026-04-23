package app

import (
	core_database "github.com/matuha07/kinotower-go/src/internal/core/database"
	core_router "github.com/matuha07/kinotower-go/src/internal/core/router"
	auth_handler "github.com/matuha07/kinotower-go/src/internal/features/auth/handler"
	auth_service "github.com/matuha07/kinotower-go/src/internal/features/auth/service"
	film_handler "github.com/matuha07/kinotower-go/src/internal/features/films/handler"
	film_repository "github.com/matuha07/kinotower-go/src/internal/features/films/repository"
	film_service "github.com/matuha07/kinotower-go/src/internal/features/films/service"
	user_handler "github.com/matuha07/kinotower-go/src/internal/features/users/handler"
	user_repository "github.com/matuha07/kinotower-go/src/internal/features/users/repository"
	user_service "github.com/matuha07/kinotower-go/src/internal/features/users/service"
)

type Container struct {
	Router *core_router.Router
}

func New(db core_database.Database, jwtSecret string) *Container {
	// Films
	filmRepo := film_repository.NewFilmRepository(db)
	filmSvc := film_service.NewFilmService(filmRepo)
	filmH := film_handler.NewFilmHandler(filmSvc)

	// Users
	userRepo := user_repository.NewUserRepository(db)
	userSvc := user_service.NewUserService(userRepo)
	userH := user_handler.NewUserHandler(userSvc)

	// Auth
	authSvc := auth_service.NewAuthService(userRepo, jwtSecret)
	authH := auth_handler.NewAuthHandler(authSvc)

	router := core_router.NewRouter(filmH, userH, authH, jwtSecret)

	return &Container{Router: router}
}
