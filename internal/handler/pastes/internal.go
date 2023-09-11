package pastes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/config"
	"github.com/merlinfuchs/vaultbin/internal/public/views"
	"github.com/merlinfuchs/vaultbin/internal/wire"
)

func (h *PastesHandler) InternalPasteCreate(c echo.Context) error {
	var req wire.InternalPasteCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	if strings.TrimSpace(req.Content) == "" || len(req.Content) > config.K.Int("paste_max_size") {
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

	ttl := time.Duration(config.K.Int("paste_default_ttl")) * time.Second
	paste, err := h.store.CreatePaste(req.Content, req.Language, ttl)
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

func (h *PastesHandler) InternalPasteDuplicate(c echo.Context) error {
	var req wire.InternalPasteDuplicateRequest
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

func (h *PastesHandler) InternalPasteNew(c echo.Context) error {
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

func (h *PastesHandler) InternalPasteRaw(c echo.Context) error {
	var req wire.InternalPasteRawRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/%s/raw", req.ID))
	return c.NoContent(http.StatusOK)
}
