package app

import (
	"html/template"
	"net/http"
	"time"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/merlinfuchs/vaultbin/internal/apierror"
	"github.com/merlinfuchs/vaultbin/internal/config"
	"github.com/merlinfuchs/vaultbin/internal/db"
	"github.com/merlinfuchs/vaultbin/internal/handler/pastes"
	"github.com/merlinfuchs/vaultbin/internal/public/static"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
	"golang.org/x/time/rate"
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

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logCtx := slog.With("method", v.Method).With("path", v.URI).With("status", v.Status).With("error", v.Error)

			if v.Error != nil {
				logCtx.Error("Request has failed")
			} else {
				logCtx.Info("Request has been processed")
			}
			return nil
		},
	}))

	if config.K.Bool("ratelimit.reverse_proxy") {
		e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	} else {
		e.IPExtractor = echo.ExtractIPDirect()
	}

	rateLimitMiddleware := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(config.K.Int("ratelimit.per_second")),
				Burst:     config.K.Int("ratelimit.burst"),
				ExpiresIn: time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	})

	pastes := pastes.New(db)

	e.Use(apierror.ErrorRewriteMiddleware())

	e.GET("/", pastes.PagePasteNew)
	e.GET("/:paste_id", pastes.PagePasteView)
	e.GET("/:paste_id/raw", pastes.PagePasteRaw)
	e.POST("/internal/pastes/create", pastes.InternalPasteCreate, rateLimitMiddleware)
	e.POST("/internal/pastes/duplicate", pastes.InternalPasteDuplicate)
	e.POST("/internal/pastes/new", pastes.InternalPasteNew)
	e.POST("/internal/pastes/raw", pastes.InternalPasteRaw)
	e.POST("/api/pastes", pastes.APIPasteCreate)
	e.GET("/api/pastes/:paste_id", pastes.APIPasteGet)

	e.StaticFS("/static", static.FS)

	return e
}
