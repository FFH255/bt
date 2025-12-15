package auth_ping_get_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type ResponseBody string

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

// Handle godoc
// @Summary      Ping
// @Description  Check if user logged in
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} ResponseBody "Pong"
// @Failure      401 {object} proto.Error "Unauthorized"
// @Router       /auth/ping [get]
// @Security     ApiKeyAuth
func (h *Handler) Handle(c *gin.Context) {
	proto.WriteJSON(c, http.StatusOK, "pong")
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/auth/ping"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}
