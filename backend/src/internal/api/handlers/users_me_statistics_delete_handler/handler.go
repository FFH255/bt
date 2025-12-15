package users_me_statistics_delete_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_statistics_delete_handler"

type statisticsService interface {
	DeleteAllForUser(ctx context.Context, userID uint64) error
}

type Handler struct {
	logger            internal.Logger
	statisticsService statisticsService
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong serverside"
	)

	switch {
	case models.IsUserNotFoundError(err):
		status = http.StatusNotFound
		message = "user not found"
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusInternalServerError:
		h.logger.Error(ctx)
	default:
		h.logger.Warning(ctx)
	}

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary Отчистить всю статистику пользователя
// @Description Делает soft delete для всей статистики пользователя
// @Tags User Statistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 "Статистика пользователя удалена"
// @Failure 401 {object} proto.Error "Пользователь не авторизован"
// @Failure 404 {object} proto.Error "Пользователь не найден"
// @Failure 500 {object} proto.Error "Внутренняя ошибка сервера (смотреть логи)"
// @Router /users/me/statistics [delete]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)
	req := newRequest(c)

	err := h.statisticsService.DeleteAllForUser(ctx, req.userID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}
}

func (h *Handler) Method() string {
	return http.MethodDelete
}

func (h *Handler) Path() string {
	return "/users/me/statistics"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(logger internal.Logger, statisticsService statisticsService) *Handler {
	return &Handler{
		logger:            logger,
		statisticsService: statisticsService,
	}
}
