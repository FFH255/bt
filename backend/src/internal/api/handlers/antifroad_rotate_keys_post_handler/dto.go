package antifroad_rotate_keys_post_handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Password string
}

func newRequest(c *gin.Context) (*Request, error) {
	password := c.Query("password")
	if password == "" {
		return nil, fmt.Errorf("password can not be empty")
	}

	return &Request{Password: password}, nil
}
