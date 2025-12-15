package users_me_username_patch_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_username_patch_handler"

type userService interface {
	ChangeUsername(ctx context.Context, userID uint64, username string) error
}

type Handler struct {
	logger      internal.Logger
	userService userService
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong serverside"
	)

	switch {
	case models.IsValidationError(err), models.IsSameUsernameError(err):
		status = http.StatusBadRequest
		message = err.Error()
	case models.IsUsernameAlreadyTakenError(err):
		status = http.StatusConflict
		message = err.Error()
	case models.IsTooMuchUsernameChangesError(err):
		status = http.StatusTooManyRequests
		message = err.Error()
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusBadRequest, http.StatusConflict, http.StatusTooManyRequests:
		h.logger.Warning(ctx)
	default:
		h.logger.Error(ctx)
	}

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary     Изменить имя пользователя
// @Description Изменить имя пользователя, но не чаще чем 1 раз в 15 дней
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       body body RequestBody true "Содержит новое имя пользователя"
// @Success     200                      "Имя пользователя обновлено"
// @Failure     400 {object} proto.Error "Имя пользователя отсутствует или не подходит под условия"
// @Failure     409 {object} proto.Error "Имя пользователя занято"
// @Failure     429 {object} proto.Error "Нельзя так часто изменять имя пользователя"
// @Failure     500 {object} proto.Error "Внутренняя ошибка сервера (смотреть логи)"
// @Router      /users/me/username [patch]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	req, err := newRequest(c)
	if err != nil {
		proto.WriteError(c, http.StatusBadRequest, err)
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		return
	}

	if err = h.userService.ChangeUsername(ctx, req.UserID, req.Username); err != nil {
		h.handleError(ctx, c, err)
		return
	}
}

func (h *Handler) Method() string {
	return http.MethodPatch
}

func (h *Handler) Path() string {
	return "users/me/username"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(logger internal.Logger, userService userService) *Handler {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}
