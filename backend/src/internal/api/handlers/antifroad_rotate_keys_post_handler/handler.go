package antifroad_rotate_keys_post_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/services/antifroad_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const (
	handlerName = "antifroad_rotate_keys_post_handler"
)

type antifroadService interface {
	RotateKeys(ctx context.Context, password string) error
}

type Handler struct {
	antifroadService antifroadService
	logger           internal.Logger
	limiter          *limiter.Limiter
}

// Handle godoc
// @Summary Ротировать ключи (служебная ручка)
// @Description Ротировать ключи. Ручка используется только внутри контура blindtyping.
// @Tags Antifroad
// @Accept  json
// @Produce  json
// @Param password query string true "Пароль для авторизации в модуле антифрода" Example(12345)
// @Success 200
// @Failure 400 {object} proto.Error "Отсутствуем пароль"
// @Failure 401 {object} proto.Error "Неправильный пароль"
// @Failure 500 {object} proto.Error "Ошибка сервера"
// @Router /antifroad/rotate-keys [post]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	limiterCtx, err := h.limiter.Get(ctx, handlerName)
	if err != nil {
		h.logger.Error(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if limiterCtx.Reached {
		h.logger.Warning(h.logger.WithStatusCode(ctx, http.StatusTooManyRequests))
		proto.WriteError(c, http.StatusTooManyRequests, "too many requests")
		return
	}

	req, err := newRequest(c)
	if err != nil {
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	err = h.antifroadService.RotateKeys(ctx, req.Password)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong"
	)

	switch {
	case antifroad_service.IsWrongPasswordError(err):
		message = err.Error()
		status = http.StatusUnauthorized
	}

	ctx = h.logger.WithStatusCode(ctx, status)
	h.logger.Error(h.logger.WithError(ctx, err))
	proto.WriteError(c, status, message)
}

func (h *Handler) Method() string {
	return http.MethodPost
}

func (h *Handler) Path() string {
	return "/antifroad/rotate-keys"
}

func (h *Handler) Middleware() []string {
	return nil
}

func New(service antifroadService, logger internal.Logger) *Handler {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  1,
	}
	store := memory.NewStore()
	limitInstance := limiter.New(store, rate)

	return &Handler{
		antifroadService: service,
		logger:           logger,
		limiter:          limitInstance,
	}
}
