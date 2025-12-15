package users_me_statistics_delete_handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal/api"
)

type Request struct {
	userID uint64
}

func newRequest(c *gin.Context) *Request {
	return &Request{
		userID: api.GetUserID(c),
	}
}
