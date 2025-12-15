package auth_refresh_post_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/services/auth_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "auth_refresh_post_handler"

type refresher interface {
	Refresh(ctx context.Context, in *auth_service.RefreshIn) (*auth_service.RefreshOut, error)
}

type cookieManager interface {
	SetAccessToken(ctx *gin.Context, accessToken string)
	SetRefreshToken(ctx *gin.Context, refreshToken string)
}

type Handler struct {
	refresher     refresher
	cookieManager cookieManager
	logger        internal.Logger
}

func New(refresher refresher, cookieManager cookieManager, logger internal.Logger) *Handler {
	return &Handler{
		refresher:     refresher,
		cookieManager: cookieManager,
		logger:        logger,
	}
}

// Handle godoc
// @Summary      Refresh
// @Description  Rotates access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} nil "Success"
// @Router       /auth/refresh [post]
// @Failure      403 {object} proto.Error "Unauthorized"
// @Failure      404 {object} proto.Error "User session not found"
// @Failure      500 {object} proto.Error "Something went wrong server side"
// @Security     ApiKeyAuth
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	out, err := h.refresher.Refresh(ctx, h.makeRefreshIn(c))
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	h.makeOut(c, out)
}

func (h *Handler) makeRefreshIn(c *gin.Context) *auth_service.RefreshIn {
	return &auth_service.RefreshIn{
		RefreshToken: api.GetRefreshToken(c),
	}
}

func (h *Handler) makeOut(c *gin.Context, out *auth_service.RefreshOut) {
	if c == nil || out == nil {
		proto.WriteError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.cookieManager.SetAccessToken(c, out.AccessToken)
	h.cookieManager.SetRefreshToken(c, out.RefreshToken)
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong"
	)

	switch {
	case auth_service.IsSessionNotFoundError(err):
		status = http.StatusNotFound
		message = "session not found"
	case auth_service.IsSessionExpiredError(err):
		status = http.StatusUnauthorized
		message = "session expired"
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusNotFound, http.StatusUnauthorized:
		h.logger.Warning(ctx)
	default:
		h.logger.Error(ctx)
	}

	proto.WriteError(c, status, message)
}

func (h *Handler) Method() string {
	return http.MethodPost
}

func (h *Handler) Path() string {
	return "/auth/refresh"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.RefreshToken}
}
