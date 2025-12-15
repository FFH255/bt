package users_me_get_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_get_handler"

type userService interface {
	GetByID(ctx context.Context, id models.ID) (*models.User, error)
}

type (
	User struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		JoinedAt string `json:"joinedAt"`
	} //@name UsersMeGetHandler.User

	ResponseBody struct {
		User User `json:"user"`
	} //@name UsersMeGetHandler.ResponseBody

	Request struct {
		ID models.ID
	}

	Handler struct {
		userService userService
		logger      internal.Logger
	}
)

func newRequest(c *gin.Context) *Request {
	return &Request{
		ID: models.ID(api.GetUserID(c)),
	}
}

func newResponseBody(user *models.User) *ResponseBody {
	return &ResponseBody{
		User: User{
			Email:    string(user.Email),
			Username: string(user.Nickname),
			JoinedAt: proto.MarshalTime(user.CreatedAt),
		},
	}
}

// Handle godoc
// @Summary Получить информацию о текущем пользователе
// @Description Получить информацию о текущем пользователе при помощи ID вшитого в access token
// @Tags Users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} ResponseBody
// @Failure 401 {object} proto.Error "Пользователь не авторизован"
// @Failure 404 {object} proto.Error "Пользователь с таким ID не найден"
// @Failure 500 {object} proto.Error "Непредвиденная ошибка сервера"
// @Router /users/me [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)
	req := newRequest(c)

	user, err := h.userService.GetByID(ctx, req.ID)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusInternalServerError)
		h.logger.Error(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	if user == nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusNotFound)
		h.logger.Warning(h.logger.WithMsg(ctx, "user not found"))
		proto.WriteError(c, http.StatusNotFound, "user not found")
		return
	}

	body := newResponseBody(user)
	proto.WriteJSON(c, http.StatusOK, body)
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/users/me"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(userService userService, logger internal.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}
