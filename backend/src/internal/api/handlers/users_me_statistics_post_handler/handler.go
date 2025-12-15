package users_me_statistics_post_handler

import (
	"context"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/statistics_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_statistics_post_handler"

type statisticsSaver interface {
	Save(ctx context.Context, in *statistics_service.SaveIn) (*statistics_service.SaveOut, error)
}

type RequestBody struct {
	WPM                        float64 `json:"wpm" example:"42.5" description:"Words per minute"`
	CPM                        float64 `json:"cpm" example:"210.3" description:"Characters per minute"`
	Accuracy                   float64 `json:"accuracy" example:"98.7" description:"Accuracy percentage"`
	DurationMs                 uint64  `json:"durationMs" example:"60000" description:"Duration in milliseconds"`
	Language                   string  `json:"language" example:"english" description:"Language used for typing"`
	Mode                       string  `json:"mode" example:"time" description:"Typing mode"`
	SubMode                    string  `json:"submode" example:"1m" description:"Mode-specific parameters"`
	IsPunctuation              bool    `json:"isPunctuation" example:"true" description:"Whether punctuation was enabled"`
	UncompletedTestsCount      *uint64 `json:"uncompletedTestsCount" example:"0" description:"Uncompleted test count"`
	UncompletedTestsDurationMs *uint64 `json:"uncompletedTestsDurationMs" example:"0" description:"Total duration of uncompleted tests"`
	UID                        string  `json:"uid" example:"0" description:"Unique request ID"`
	Sign                       string  `json:"sign" example:"12345" description:"Signature"`
	CreatedAt                  string  `json:"createdAt" example:"2025-10-19T19:02:29+03:00" description:"Creation time in RFC3339"`
	StartedAt                  string  `json:"startedAt" example:"2025-10-19T19:02:29+03:00" description:"Start time in RFC3339"`
	FinishedAt                 string  `json:"finishedAt" example:"2025-10-19T19:02:29+03:00" description:"Finish time in RFC3339"`
} //@name UsersMeStatisticsPostHandler.RequestBody

type Request struct {
	body   *RequestBody
	userID models.ID
}

type ResponseBody struct {
	IsPersonalBest bool    `json:"isPersonalBest"`
	WPMShift       float64 `json:"wpmShift"`
} //@name UsersMyStatisticsPostHandler.ResponseBody

type Handler struct {
	statisticsSaver statisticsSaver
	logger          internal.Logger
}

func (h *Handler) newRequest(c *gin.Context) (*Request, error) {
	body := new(RequestBody)
	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		return nil, err
	}

	return &Request{
		body:   body,
		userID: models.ID(api.GetUserID(c)),
	}, nil
}

func (h *Handler) newSaveIn(r *Request) (*statistics_service.SaveIn, error) {
	createdAt, err := proto.UnmarshalTime(r.body.CreatedAt)
	if err != nil {
		return nil, err
	}

	startedAt, err := proto.UnmarshalTime(r.body.StartedAt)
	if err != nil {
		return nil, err
	}

	finishedAt, err := proto.UnmarshalTime(r.body.FinishedAt)
	if err != nil {
		return nil, err
	}

	return &statistics_service.SaveIn{
		UserID:                     r.userID,
		WPM:                        r.body.WPM,
		CPM:                        r.body.CPM,
		Accuracy:                   r.body.Accuracy,
		Duration:                   proto.ParseMilliseconds(&r.body.DurationMs),
		Language:                   r.body.Language,
		Mode:                       r.body.Mode,
		SubMode:                    r.body.SubMode,
		IsPunctuation:              r.body.IsPunctuation,
		UncompletedTestsCount:      pointer.GetUint64(r.body.UncompletedTestsCount),
		UncompletedTestsDurationMs: proto.ParseMilliseconds(r.body.UncompletedTestsDurationMs),
		UID:                        r.body.UID,
		Sign:                       r.body.Sign,
		CreatedAt:                  createdAt,
		StartedAt:                  startedAt,
		FinishedAt:                 finishedAt,
	}, nil
}

func (h *Handler) newResponseBody(out *statistics_service.SaveOut) *ResponseBody {
	return &ResponseBody{
		IsPersonalBest: out.IsPB,
		WPMShift:       out.WPMShift,
	}
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong serverside"
	)

	switch {
	case statistics_service.IsUserNotFoundError(err):
		status = http.StatusNotFound
		message = "user not found"
	case models.IsValidationError(err):
		status = http.StatusBadRequest
		message = err.Error()
	case statistics_service.IsFroadError(err):
		status = http.StatusBadRequest
		message = "froad detected"
	case statistics_service.IsAlreadyHandledError(err):
		status = http.StatusConflict
		message = "stats already handled"
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
// @Summary Save typing statistics
// @Description Creates a new statistics record for the authenticated user
// @Tags User Statistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body RequestBody true "Statistics data to save"
// @Success 201 {object} ResponseBody "Statistics successfully created"
// @Failure 400 {object} proto.Error "Invalid request body"
// @Failure 401 {object} proto.Error "Unauthorized"
// @Failure 500 {object} proto.Error "Internal server error"
// @Router /users/me/statistics [post]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	r, err := h.newRequest(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	in, err := h.newSaveIn(r)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	saveOut, err := h.statisticsSaver.Save(ctx, in)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	responseBody := h.newResponseBody(saveOut)
	proto.WriteJSON(c, http.StatusCreated, responseBody)
}

func (h *Handler) Method() string {
	return http.MethodPost
}

func (h *Handler) Path() string {
	return "/users/me/statistics"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(statisticsSaver statisticsSaver, logger internal.Logger) *Handler {
	return &Handler{
		statisticsSaver: statisticsSaver,
		logger:          logger,
	}
}
