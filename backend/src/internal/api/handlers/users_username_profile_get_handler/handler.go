package users_username_profile_get_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/profile_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_username_profile_get_handler"

type profileGetter interface {
	Get(ctx context.Context, in *profile_service.GetIn) (*models.Profile, error)
}

type Handler struct {
	profileGetter profileGetter
	logger        internal.Logger
}

func (h *Handler) newProfileServiceGetIn(request *Request) *profile_service.GetIn {
	return &profile_service.GetIn{
		Username: request.Username,
	}
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong"
	)

	switch {
	case profile_service.IsUserNotFoundError(err):
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
// @Summary Gets user's profile
// @Description Gets profile of any user by username
// @Tags Profile
// @Accept json
// @Produce json
// @Param username path string true "Username of user for profile fetching"
// @Success 200 {object} ResponseBody "User's profile"
// @Failure 400 {object} proto.Error "Invalid request body"
// @Failure 401 {object} proto.Error "Unauthorized"
// @Failure 404 {object} proto.Error "User not found"
// @Failure 500 {object} proto.Error "Internal server error"
// @Router /users/{username}/profile [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	request, err := newRequest(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	profile, err := h.profileGetter.Get(ctx, h.newProfileServiceGetIn(request))
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	body := newResponseBody(profile)
	c.JSON(http.StatusOK, body)
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/users/:username/profile"
}

func (h *Handler) Middleware() []string {
	return nil
}

func New(profileGetter profileGetter, logger internal.Logger) *Handler {
	return &Handler{
		profileGetter: profileGetter,
		logger:        logger,
	}
}
