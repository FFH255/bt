package antifroad_key_get_handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Key string `json:"key" example:"12345" description:"Антифрод-ключ"`
} //@name AntifroadKeyGetHandler.ResponseBody

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
