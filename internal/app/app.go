package app

import (
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/merlinfuchs/vaultbin/internal/db"
	"github.com/merlinfuchs/vaultbin/internal/handler/pastes"
	"github.com/merlinfuchs/vaultbin/internal/public/static"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

type Template struct {
	templates *template.Template
}

func New(db *db.DB) *echo.Echo {
	e := echo.New()

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 6,
	}))

	t := views.New()
	e.Renderer = t

	pastes := pastes.New(db)

	e.GET("/", pastes.PasteNew).Name = "paste_new"
	e.GET("/:paste_id", pastes.PasteView).Name = "paste_view"
	e.GET("/:paste_id/raw", pastes.PasteRaw).Name = "paste_raw"
	e.POST("/pastes/create", pastes.PasteAPICreate).Name = "paste_create"
	e.POST("/pastes/duplicate", pastes.PasteAPIDuplicate).Name = "paste_duplicate"
	e.POST("/pastes/new", pastes.PasteAPINew).Name = "paste_new"
	e.POST("/pastes/raw", pastes.PasteAPIRaw).Name = "paste_new"

	e.StaticFS("/static", static.FS)

	return e
}
