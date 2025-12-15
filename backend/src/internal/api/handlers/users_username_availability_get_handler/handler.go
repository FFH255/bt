package users_username_availability_get_handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_username_availability_get_handler"

type (
	Request struct {
		Username string
	}

	ResponseBody struct {
		Available bool `json:"available"`
	} //@name UsersUsernameAvailabilityGetHandler.ResponseBody
)

type availabilityService interface {
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
}

type Handler struct {
	availabilityService availabilityService
	logger              internal.Logger
}

func (h *Handler) newRequest(c *gin.Context) (*Request, error) {
	username := c.Query("username")
	if username == "" {
		return nil, errors.New("username is required")
	}

	return &Request{Username: username}, nil
}

func (h *Handler) newResponseBody(isAvailable bool) *ResponseBody {
	return &ResponseBody{Available: isAvailable}
}

// Handle godoc
// @Summary      Проверяет username на доступность для регистрации
// @Description  Если username свободен и соответствует требованиям, то в теле ответа будет available: true, иначе available: false. В cookie должен быть живой токен регистрации (registration_token), иначе 401
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        username query string true "Username, который нужно проверить на доступность" Example(ffh)
// @Success      200 {object} ResponseBody
// @Failure      400 {object} proto.Error
// @Failure      500 {object} proto.Error
// @Router       /users/username-availability [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	req, err := h.newRequest(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	isAvailable, err := h.availabilityService.IsUsernameAvailable(ctx, req.Username)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusInternalServerError)
		h.logger.Error(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	body := h.newResponseBody(isAvailable)
	c.JSON(http.StatusOK, body)
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/users/username-availability"
}

func (h *Handler) Middleware() []string {
	return nil
}

func New(availabilityService availabilityService, logger internal.Logger) *Handler {
	return &Handler{
		availabilityService: availabilityService,
		logger:              logger,
	}
}
