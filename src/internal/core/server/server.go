package core_server

import (
	"net/http"

	"github.com/matuha07/kinotower-go/src/internal/app"
	core_database "github.com/matuha07/kinotower-go/src/internal/core/database"
)

type Server struct {
	http.Server
}

func NewServer(db core_database.Database) *Server {
	cfg := NewConfigMust()
	container := app.New(db, cfg.JWTSecret)

	return &Server{
		Server: http.Server{
			Addr:    cfg.Addr,
			Handler: container.Router,
		},
	}
}
