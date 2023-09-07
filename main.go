package main

import (
	"log/slog"

	"github.com/merlinfuchs/vaultbin/internal/app"
	"github.com/merlinfuchs/vaultbin/internal/db"
)

func main() {
	db, err := db.New()
	if err != nil {
		slog.With("error", err).Error("Error opening database")
		return
	}

	e := app.New(db)
	err = e.Start(":8080")
	if err != nil {
		slog.With("error", err).Error("Error starting server")
	}
}
