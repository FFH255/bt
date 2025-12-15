package auth_provider_callback_post_handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/services/auth_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/oauth"
)

type Status int

const (
	EmailIsAlreadyTaken Status = 1

	statusParamName = "status"

	handlerName = "auth_provider_callback_post_handler"
)

type authService interface {
	Callback(ctx context.Context, in *auth_service.CallbackIn) (*auth_service.CallbackOut, error)
}

type cookieManager interface {
	SetAccessToken(c *gin.Context, accessToken string)
	SetRefreshToken(c *gin.Context, refreshToken string)
	SetRegistrationToken(c *gin.Context, registrationToken string)
}

type oauthManager interface {
	Complete(c *gin.Context, provider string) (*oauth.User, error)
}

type Handler struct {
	oauthManager    oauthManager
	cookieManager   cookieManager
	authService     authService
	logger          internal.Logger
	loggedInURL     string
	registrationURL string
	errorURL        string
}

func New(
	authService authService,
	oauthManager oauthManager,
	cookieManager cookieManager,
	logger internal.Logger,
	loggedInURL string,
	registrationURL string,
	errorURL string,
) *Handler {
	return &Handler{
		authService:     authService,
		oauthManager:    oauthManager,
		cookieManager:   cookieManager,
		logger:          logger,
		loggedInURL:     loggedInURL,
		registrationURL: registrationURL,
		errorURL:        errorURL,
	}
}

func (h *Handler) buildErrorURL(status Status) string {
	u, err := url.Parse(h.errorURL)
	if err != nil {
		panic(err)
	}

	params := url.Values{}
	params.Add(statusParamName, strconv.Itoa(int(status)))

	u.RawQuery = params.Encode()

	return u.String()
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case auth_service.IsEmailAlreadyTakenError(err):
		c.Redirect(http.StatusPermanentRedirect, h.buildErrorURL(EmailIsAlreadyTaken))
	default:
		c.Redirect(http.StatusPermanentRedirect, h.errorURL)
	}
}

// Handle godoc
// @Summary      OAuth Callback Handler
// @Description  Обрабатывает callback-запрос от OAuth провайдера. Завершает процесс аутентификации и выполняет перенаправление на соответствующий фронтенд URL в зависимости от результата.
// @Description
// @Description  **Сценарии перенаправления:**
// @Description  - Успешная аутентификация существующего пользователя → `logged_in_url`
// @Description  - Первая аутентификация нового пользователя → `registration_url`
// @Description  - Ошибка аутентификации → `error_redirection_url?status=STATUS`
// @Description
// @Description  **Коды ошибок:**
// @Description  - `1` - Email уже занят другие аккаунтом
// @Description
// @Description  **Особенности:**
// @Description  - Устанавливает authentication cookies при успешной аутентификации или registration_token при необходимости регистрации
// @Description  - Всегда возвращает HTTP 308 (Permanent Redirect)
// @Description  - Параметр `status` передается только в случае ошибки
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider path string true "Название OAuth провайдера (google или github)" enum(google,github)
// @Success      308 "Перенаправление на фронтенд"
// @Router       /auth/{provider}/callback [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	in, err := newIn(c)
	if err != nil {
		h.logger.Error(h.logger.WithError(ctx, err))
		c.Redirect(http.StatusPermanentRedirect, h.errorURL)
		return
	}

	ctx = h.logger.WithField(ctx, "provider", in.Provider)

	user, err := h.oauthManager.Complete(c, in.Provider)
	if err != nil {
		h.logger.Error(h.logger.WithError(ctx, err))
		c.Redirect(http.StatusPermanentRedirect, h.errorURL)
		return
	}
	if user == nil {
		h.logger.Error(h.logger.WithField(ctx, "msg", "oauth user is nil"))
		c.Redirect(http.StatusPermanentRedirect, h.errorURL)
		return
	}

	ctx = h.logger.WithField(ctx, "external_user_id", user.ID)

	out, err := h.authService.Callback(ctx, &auth_service.CallbackIn{
		Email:    user.Email,
		Provider: user.Provider,
	})
	if err != nil {
		h.logger.Error(h.logger.WithError(ctx, err))
		h.handleError(c, err)
		return
	}
	if out == nil {
		h.logger.Error(h.logger.WithField(ctx, "msg", "auth service response is nil"))
		c.Redirect(http.StatusPermanentRedirect, h.errorURL)
		return
	}

	if out.IsRegistration() {
		h.cookieManager.SetRegistrationToken(c, *out.RegistrationToken)
		c.Redirect(http.StatusPermanentRedirect, h.registrationURL)
		return
	}

	if out.IsLogin() {
		h.cookieManager.SetAccessToken(c, *out.AccessToken)
		h.cookieManager.SetRefreshToken(c, *out.RefreshToken)
		c.Redirect(http.StatusPermanentRedirect, h.loggedInURL)
		return
	}

	h.logger.Error(h.logger.WithField(ctx, "msg", "out is not registration and login"))
	c.Redirect(http.StatusPermanentRedirect, h.errorURL)
}

func (*Handler) Method() string {
	return http.MethodGet
}

func (*Handler) Path() string {
	return "/auth/:provider/callback"
}

func (h *Handler) Middleware() []string {
	return nil
}
