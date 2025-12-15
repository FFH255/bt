package auth_logout_post_handler

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

const (
	handlerName = "auth_logout_post_handler"
)

type authService interface {
	Logout(ctx context.Context, in *auth_service.LogoutIn) error
}

type oauth interface {
	Logout(c *gin.Context) error
}

type cookieManager interface {
	DeleteAccessToken(ctx *gin.Context)
	DeleteRefreshToken(ctx *gin.Context)
}

type Handler struct {
	authService   authService
	cookieManager cookieManager
	oauth         oauth
	logger        internal.Logger
}

func New(
	authService authService,
	cookieManager cookieManager,
	oauth oauth,
	logger internal.Logger,
) *Handler {
	return &Handler{
		authService:   authService,
		cookieManager: cookieManager,
		oauth:         oauth,
		logger:        logger,
	}
}

// Handle godoc
// @Summary      Logout
// @Description  Delete session (refresh token), add access token to block list for some time and delete tokens from cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 "Success"
// @Failure      403 {object} proto.Error "Unauthorized"
// @Failure 	 404 {object} proto.Error "User session not found"
// @Failure		 500 {object} proto.Error "Server side error"
// @Router       /auth/logout [post]
// @Security     ApiKeyAuth
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	err := h.oauth.Logout(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusInternalServerError)
		h.logger.Error(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	err = h.authService.Logout(c, h.makeLogoutIn(c))
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	h.cookieManager.DeleteAccessToken(c)
	h.cookieManager.DeleteRefreshToken(c)
}

func (h *Handler) makeLogoutIn(c *gin.Context) *auth_service.LogoutIn {
	return &auth_service.LogoutIn{
		AccessToken:  api.GetAccessToken(c),
		RefreshToken: api.GetRefreshToken(c),
	}
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
	return "/auth/logout"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}
