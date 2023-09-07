package pastes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

type PasteAPICreateRequest struct {
	Content  string `form:"content"`
	Language string `form:"language"`
}

func (h *PastesHandler) PasteAPICreate(c echo.Context) error {
	var req PasteAPICreateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	paste, err := h.store.CreatePaste(req.Content, req.Language)
	if err != nil {
		slog.With("error", err).Error("Failed to create paste")
		return err
	}

	_, err = h.store.CountPasteView(paste.ID)
	if err != nil {
		slog.With("error", err).Error("Failed to count initial paste view")
		return err
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/%s", paste.ID))
	err = c.Render(http.StatusOK, "paste.content", views.PasteViewData{
		New:      false,
		PasteID:  paste.ID,
		Content:  paste.Content,
		Language: paste.Language,
		Views:    1,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}
