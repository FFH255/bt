package auth_provider_get_handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type in struct {
	Provider string `json:"provider"`
}

func newIn(c *gin.Context) (*in, error) {
	provider := c.Param("provider")
	if provider == "" {
		return nil, errors.New("provider is required")
	}

	return &in{Provider: provider}, nil
}
