package pastes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *PastesHandler) PasteRaw(c echo.Context) error {
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
