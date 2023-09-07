package pastes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PasteAPIRawRequest struct {
	ID string `form:"id"`
}

func (h *PastesHandler) PasteAPIRaw(c echo.Context) error {
	var req PasteAPIRawRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/%s/raw", req.ID))
	return c.NoContent(http.StatusOK)
}
