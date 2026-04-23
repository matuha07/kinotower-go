package main

import (
	"log/slog"

	core_database "github.com/matuha07/kinotower-go/src/internal/core/database"
	core_logger "github.com/matuha07/kinotower-go/src/internal/core/logger"
	core_server "github.com/matuha07/kinotower-go/src/internal/core/server"
)

func main() {
	logger := core_logger.FromEnv("kinotower")
	slog.SetDefault(logger)
	logger.Info("starting server", "addr", ":8080")

	db, err := core_database.NewDatabase()
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return
	}
	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		return
	}

	server := core_server.NewServer(*db)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped with error", "error", err)
	}
}
