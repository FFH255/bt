package auth_provider_get_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "auth_provider_get_handler"

type oauthHandler interface {
	Begin(c *gin.Context, provider string)
}

type Handler struct {
	oauth  oauthHandler
	logger internal.Logger
}

func New(
	oauth oauthHandler,
	logger internal.Logger,
) *Handler {
	return &Handler{
		oauth:  oauth,
		logger: logger,
	}
}

// Handle godoc
// @Summary      OAuth Login Redirect
// @Description  Initiates the OAuth authentication process by redirecting the user to the provider's login page.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider path string true "OAuth provider name"
// @Success      200 {object} nil "User already authenticated"
// @Failure      400 {object} proto.Error "Missing provider or invalid input"
// @Router       /auth/{provider} [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	in, err := newIn(c)
	if err != nil {
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	h.oauth.Begin(c, in.Provider)
}

func (*Handler) Method() string {
	return http.MethodGet
}

func (*Handler) Path() string {
	return "/auth/:provider"
}

func (h *Handler) Middleware() []string {
	return nil
}
