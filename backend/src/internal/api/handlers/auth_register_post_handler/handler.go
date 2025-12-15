package auth_register_post_handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal"

	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/auth_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/user_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "auth_register_post_handler"

type RequestBody struct {
	Nickname string `json:"nickname" example:"ffh255"`
} //@name AuthRegisterPostHandler.RequestBody

type Request struct {
	Nickname string
	Token    string
}

type authService interface {
	Register(ctx context.Context, in *auth_service.RegisterIn) (*auth_service.RegisterOut, error)
}

type cookieManager interface {
	DeleteRegistrationToken(c *gin.Context)
	SetAccessToken(c *gin.Context, token string)
	SetRefreshToken(c *gin.Context, token string)
}

type Handler struct {
	authService   authService
	cookieManager cookieManager
	logger        internal.Logger
}

func New(authService authService, cookieManager cookieManager, logger internal.Logger) *Handler {
	return &Handler{
		authService:   authService,
		cookieManager: cookieManager,
		logger:        logger,
	}
}

func (h *Handler) newRequest(c *gin.Context) (*Request, error) {
	var body RequestBody
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		return nil, errors.New("failed to parse json body")
	}

	return &Request{
		Nickname: body.Nickname,
		Token:    api.GetRegistrationToken(c),
	}, nil
}

func (h *Handler) writeResponse(c *gin.Context, out *auth_service.RegisterOut) {
	h.cookieManager.DeleteRegistrationToken(c)
	h.cookieManager.SetAccessToken(c, out.AccessToken)
	h.cookieManager.SetRefreshToken(c, out.RefreshToken)
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong serverside"
	)

	switch {
	case models.IsValidationError(err):
		status = http.StatusBadRequest
		message = err.Error()
	case user_service.IsUserAlreadyExistsError(err):
		status = http.StatusConflict
		message = "user already exists"
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusBadRequest, http.StatusConflict:
		h.logger.Warning(ctx)
	default:
		h.logger.Error(ctx)
	}

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary      Register user
// @Description  Finalize user registration by submitting a nickname and a registration token (via cookie).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RequestBody true "User nickname"
// @Success      200 "Access and refresh tokens are set as cookies"
// @Failure      400   {object}  proto.Error  "Invalid input or missing registration token"
// @Failure      409   {object}  proto.Error  "User already exists"
// @Failure      500   {object}  proto.Error  "Internal server error during registration"
// @Router       /auth/register [post]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	r, err := h.newRequest(c)
	if err != nil {
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	out, err := h.authService.Register(ctx, &auth_service.RegisterIn{
		Nickname: r.Nickname,
		Token:    r.Token,
	})
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	h.writeResponse(c, out)
}

func (h *Handler) Path() string {
	return "/auth/register"
}

func (h *Handler) Method() string {
	return http.MethodPost
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Registration}
}
