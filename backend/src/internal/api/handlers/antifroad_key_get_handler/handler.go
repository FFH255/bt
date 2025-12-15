package antifroad_key_get_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/antifroad_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const (
	handlerName = "antifroad_key_get_handler"
)

type antifroadService interface {
	GetKey(password string) (*models.AntifroadKey, error)
}

type Handler struct {
	antifroadService antifroadService
	logger           internal.Logger
	limiter          *limiter.Limiter
}

// Handle godoc
// @Summary Получить актуальный антофрод-ключ
// @Description Получить актуальный антофрод-ключ. Ручка используется только внутри контура blindtyping.
// @Tags Antifroad
// @Accept  json
// @Produce  json
// @Param password query string true "Пароль для авторизации в модуле антифрода" Example(12345)
// @Success 200 {object} ResponseBody "Антифрод-ключ"
// @Failure 400 {object} proto.Error "Отсутствуем пароль"
// @Failure 401 {object} proto.Error "Неправильный пароль"
// @Failure 404 {object} proto.Error "Ключ не найден (вообще такого быть не должно)"
// @Router /antifroad/key [get]
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

	key, err := h.antifroadService.GetKey(req.Password)
	if err != nil {
		h.handleError(c, ctx, err)
		return
	}
	if key == nil || key.Value == "" {
		h.logger.Error(h.logger.WithMsg(ctx, "antifroad key not found"))
		proto.WriteError(c, http.StatusNotFound, "key not found")
		return
	}

	out := h.newResponseBody(key)

	proto.WriteJSON(c, http.StatusOK, out)
}

func (h *Handler) handleError(c *gin.Context, ctx context.Context, err error) {
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

func (h *Handler) newResponseBody(key *models.AntifroadKey) *ResponseBody {
	return &ResponseBody{
		Key: key.Value,
	}
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/antifroad/key"
}

func (h *Handler) Middleware() []string {
	return nil
}

func New(antifroadService antifroadService, logger internal.Logger) *Handler {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  1,
	}
	store := memory.NewStore()
	limitInstance := limiter.New(store, rate)

	return &Handler{
		antifroadService: antifroadService,
		logger:           logger,
		limiter:          limitInstance,
	}
}
