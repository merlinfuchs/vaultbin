package main

import (
	"fmt"
	"log/slog"

	"github.com/merlinfuchs/vaultbin/internal/app"
	"github.com/merlinfuchs/vaultbin/internal/config"
	"github.com/merlinfuchs/vaultbin/internal/db"
)

func main() {
	config.InitConfig()

	db, err := db.New()
	if err != nil {
		slog.With("error", err).Error("Error opening database")
		return
	}

	e := app.New(db)
	err = e.Start(fmt.Sprintf("%s:%s", config.K.String("host"), config.K.String("port")))
	if err != nil {
		slog.With("error", err).Error("Error starting server")
	}
}
