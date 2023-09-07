package pastes

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

type PasteAPIDuplicateRequest struct {
	Content string `form:"content"`
}

func (h *PastesHandler) PasteAPIDuplicate(c echo.Context) error {
	var req PasteAPIDuplicateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Push-Url", "/")
	err = c.Render(http.StatusOK, "paste.content", views.PasteViewData{
		New:     true,
		Content: req.Content,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}
