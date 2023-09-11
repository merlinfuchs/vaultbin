package pastes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/merlinfuchs/vaultbin/internal/apierror"
	"github.com/merlinfuchs/vaultbin/internal/config"
	"github.com/merlinfuchs/vaultbin/internal/wire"
)

func (h *PastesHandler) APIPasteGet(c echo.Context) error {
	pasteID := c.Param("paste_id")

	paste, err := h.store.Paste(pasteID)
	if err != nil {
		return err
	}

	if paste == nil {
		return apierror.Error(apierror.ErrorCodeNotFound, "paste doesn't exist or has expired")
	}

	viewCount, err := h.store.CountPasteView(paste.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, wire.APIPasteGetResponse{
		Success: true,
		Data: wire.APIPasteGetResponseData{
			ID:        paste.ID,
			Content:   paste.Content,
			Language:  paste.Language,
			CreatedAt: paste.CreatedAt,
			ExpiresAt: paste.ExpiresAt,
			ViewCount: viewCount,
		},
	})
}

func (h *PastesHandler) APIPasteCreate(c echo.Context) error {
	var req wire.APIPasteCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	pasteMaxSize := config.K.Int("paste_max_size")
	if strings.TrimSpace(req.Content) == "" || len(req.Content) > pasteMaxSize {
		return apierror.Error(apierror.ErrorCodeBadRequest, fmt.Sprintf("paste content must be between 1 and %d characters", pasteMaxSize))
	}

	maxTTL := time.Duration(config.K.Int("paste_max_ttl")) * time.Second
	ttl := time.Duration(req.Expiration) * time.Second
	if ttl > maxTTL {
		return apierror.Error(apierror.ErrorCodeBadRequest, fmt.Sprintf("paste expiration must be between 1 and %d seconds", int(maxTTL.Seconds())))
	} else if ttl == 0 {
		ttl = time.Duration(config.K.Int("paste_default_ttl")) * time.Second
	}

	paste, err := h.store.CreatePaste(req.Content, req.Language, ttl)
	if err != nil {
		return fmt.Errorf("failed to create paste: %w", err)
	}

	viewCount, err := h.store.CountPasteView(paste.ID)
	if err != nil {
		return fmt.Errorf("failed to count initial paste view: %w", err)
	}

	return c.JSON(http.StatusOK, wire.APIPasteCreateResponse{
		Success: true,
		Data: wire.APIPasteCreateResponseData{
			ID:        paste.ID,
			Content:   paste.Content,
			Language:  paste.Language,
			CreatedAt: paste.CreatedAt,
			ExpiresAt: paste.ExpiresAt,
			ViewCount: viewCount,
		},
	})
}
