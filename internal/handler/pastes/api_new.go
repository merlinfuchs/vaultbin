package pastes

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

func (h *PastesHandler) PasteAPINew(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/")
	err := c.Render(http.StatusOK, "paste.content", views.PasteViewData{
		New: true,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}
