package pastes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

func (h *PastesHandler) PasteView(c echo.Context) error {
	pasteID := c.Param("paste_id")

	paste, err := h.store.Paste(pasteID)
	if err != nil {
		return err
	}
	if paste == nil {
		return fmt.Errorf("Paste not found")
	}

	viewCount, err := h.store.CountPasteView(paste.ID)
	if err != nil {
		return err
	}

	err = c.Render(http.StatusOK, "paste", views.PasteViewData{
		New:      false,
		PasteID:  paste.ID,
		Content:  paste.Content,
		Language: paste.Language,
		Views:    viewCount,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}
