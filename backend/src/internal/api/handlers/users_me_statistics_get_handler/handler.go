package users_me_statistics_get_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/statistics_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_statistics_get_handler"

type Request struct {
	UserID   models.ID
	DateFrom *time.Time
}

type Stats struct {
	ID                            uint    `json:"id" example:"1" format:"uint"`
	WPM                           float64 `json:"wpm" example:"42.5" format:"float64"`
	CPM                           float64 `json:"cpm" example:"210.3" format:"float64"`
	Accuracy                      float64 `json:"accuracy" example:"98.7" format:"float64"`
	Duration                      int64   `json:"duration" example:"120000" format:"int64" description:"Duration in microseconds"`
	PlayedAt                      string  `json:"playedAt" example:"2023-12-25T15:04:05Z" format:"date-time"`
	Language                      string  `json:"language" example:"english"`
	Mode                          string  `json:"mode" example:"time"`
	SubMode                       string  `json:"submode" example:"30s"`
	IsPunctuation                 bool    `json:"isPunctuation" example:"true"`
	UncompletedTestsCount         int64   `json:"uncompletedTestsCount" example:"0"`
	UncompletedTestsTotalDuration int64   `json:"uncompletedTestsTotalDuration" example:"0"`
} //@name UsersMeStatisticsGetHandler.Stats

type ResponseBody struct {
	Statistics []Stats `json:"statistics"`
} //@name UsersMeStatisticsGetHandler.ResponseBody

type statisticsGetter interface {
	GetByUser(ctx context.Context, in *statistics_service.GetByUserIn) ([]models.Statistics, error)
}

type Handler struct {
	statisticsGetter statisticsGetter
	logger           internal.Logger
}

func (h *Handler) newRequest(c *gin.Context) (*Request, error) {
	r := &Request{
		UserID: models.ID(api.GetUserID(c)),
	}

	dateFromString := c.Query("dateFrom")
	if dateFromString == "" {
		return r, nil
	}

	dateFrom, err := proto.UnmarshalTime(dateFromString)
	if err != nil {
		return nil, err
	}

	r.DateFrom = &dateFrom

	return r, nil
}

func (h *Handler) newStatisticsGetByUser(r *Request) *statistics_service.GetByUserIn {
	return &statistics_service.GetByUserIn{
		UserID:   r.UserID,
		DateFrom: r.DateFrom,
	}
}

func (h *Handler) newResponseBody(stats []models.Statistics) ResponseBody {
	body := make([]Stats, len(stats))

	for i, stat := range stats {
		body[i] = Stats{
			ID:                            uint(stat.ID),
			WPM:                           float64(stat.WPM),
			CPM:                           float64(stat.CPM),
			Accuracy:                      float64(stat.Accuracy),
			Duration:                      stat.Duration.Milliseconds(),
			PlayedAt:                      proto.MarshalTime(stat.PlayedAt),
			Language:                      string(stat.Language),
			Mode:                          string(stat.Mode),
			SubMode:                       string(stat.SubMode),
			IsPunctuation:                 stat.IsPunctuation,
			UncompletedTestsCount:         int64(stat.UncompletedTestsCount),
			UncompletedTestsTotalDuration: stat.UncompletedTestsTotalDuration.Milliseconds(),
		}
	}

	return ResponseBody{
		Statistics: body,
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
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusNotFound:
		h.logger.Warning(ctx)
	default:
		h.logger.Error(ctx)
	}

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary Получить статистику текущего пользователя
// @Description Отдает статистику текущего пользователя
// @Tags User Statistics
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param dateFrom query string false "Дата начала периода выборки в RFC3339" Example(2015-09-15T14:00:12-00:00)
// @Success 200 {object} ResponseBody "Статистика пользователя"
// @Failure 400 {object} proto.Error "Неправильно сформирован запрос"
// @Failure 401 {object} proto.Error "Пользователь не авторизован в системе"
// @Failure 404 {object} proto.Error "Пользователь не найден"
// @Failure 500 {object} proto.Error "Серверная ошибка"
// @Router /users/me/statistics [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	r, err := h.newRequest(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	stats, err := h.statisticsGetter.GetByUser(ctx, h.newStatisticsGetByUser(r))
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	out := h.newResponseBody(stats)
	proto.WriteJSON(c, http.StatusOK, out)
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/users/me/statistics"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(statisticsGetter statisticsGetter, logger internal.Logger) *Handler {
	return &Handler{
		statisticsGetter: statisticsGetter,
		logger:           logger,
	}
}
