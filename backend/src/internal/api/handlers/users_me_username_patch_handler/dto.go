package users_me_username_patch_handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal/api"
)

type RequestBody struct {
	Username string `json:"username"`
} //@name UserUsernamePatchHandler.RequestBody

type Request struct {
	Username string
	UserID   uint64
}

func newRequest(c *gin.Context) (*Request, error) {
	var body RequestBody

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		return nil, err
	}
	userID := api.GetUserID(c)

	return &Request{
		Username: body.Username,
		UserID:   userID,
	}, nil
}
