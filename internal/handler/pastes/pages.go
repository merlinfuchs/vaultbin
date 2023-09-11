package pastes

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
)

func (h *PastesHandler) PagePasteNew(c echo.Context) error {
	err := c.Render(http.StatusOK, "paste", views.PasteViewData{
		New: true,
	})
	if err != nil {
		slog.With("error", err).Error("failed to render paste_new template")
		return err
	}
	return nil
}

func (h *PastesHandler) PagePasteRaw(c echo.Context) error {
	pasteID := c.Param("paste_id")

	paste, err := h.store.Paste(pasteID)
	if err != nil {
		return err
	}

	if paste == nil {
		return c.String(http.StatusOK, "Paste doesn't exist or has expired")
	}

	return c.String(http.StatusOK, paste.Content)
}

func (h *PastesHandler) PagePasteView(c echo.Context) error {
	pasteID := c.Param("paste_id")

	paste, err := h.store.Paste(pasteID)
	if err != nil {
		return err
	}

	if paste == nil {
		err = c.Render(http.StatusOK, "paste", views.PasteViewData{
			New:     true,
			Content: "Paste doesn't exist or has expired",
		})
		if err != nil {
			slog.With("error", err).Error("failed to render paste_new template")
			return err
		}
		return nil
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
