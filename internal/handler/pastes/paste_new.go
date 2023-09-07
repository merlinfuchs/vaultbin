package pastes

import (
	"net/http"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

func (h *PastesHandler) PasteNew(c echo.Context) error {
	err := c.Render(http.StatusOK, "paste", views.PasteViewData{
		New: true,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}
