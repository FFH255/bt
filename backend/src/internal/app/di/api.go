package di

import (
	"github.com/gin-contrib/cors"

	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/antifroad_key_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/antifroad_rotate_keys_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_logout_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_ping_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_provider_callback_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_provider_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_refresh_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/auth_register_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_me_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_me_statistics_delete_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_me_statistics_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_me_statistics_post_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_me_username_patch_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_username_availability_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/handlers/users_username_profile_get_handler"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware/auth_middleware"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware/refresh_token_middleware"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware/registration_middleware"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware/request_id_middleware"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto/swagger"
)

func (c *Container) Router() *proto.Router {
	if c.router == nil {
		corsConfig := c.cfg.CORS
		swaggerConfig := c.cfg.Swagger

		router := proto.NewRouter()

		// CORS
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  corsConfig.AllowAllOrigins,
			AllowOrigins:     corsConfig.AllowedOrigins,
			AllowMethods:     corsConfig.AllowedMethods,
			AllowHeaders:     corsConfig.AllowedHeaders,
			AllowCredentials: corsConfig.AllowCredentials,
			ExposeHeaders:    corsConfig.ExposedHeaders,
		}))

		// RequestID
		router.Use(request_id_middleware.New(c.Logger(), c.UUIDGenerator()))

		// Middlewares
		router.Middleware(
			c.AuthMiddleware(),
			c.RegistrationMiddleware(),
			c.RefreshTokenMiddleware(),
		)

		// Handlers
		router.Handle(
			swagger.New(swaggerConfig.Login, swaggerConfig.Password),
			c.AuthLogoutPostHandler(),
			c.AuthPingGetHandler(),
			c.AuthProviderCallbackGetHandler(),
			c.AuthProviderGetHandler(),
			c.AuthRefreshPostHandler(),
			c.AuthRegisterPostHandler(),
			c.UsersUsernameAvailabilityGetHandler(),
			c.UsersMeStatisticsGetHandler(),
			c.UsersMeStatisticsPostHandler(),
			c.UsersMeStatisticsDeleteHandler(),
			c.UsersUsernameProfileGetHandler(),
			c.UsersMeGetHandler(),
			c.AntifroadKeyGetHandler(),
			c.AntifroadRotateKeysPostHandler(),
			c.UsersMeUsernamePatchHandler(),
		)

		c.router = router
	}
	return c.router
}

func (c *Container) AuthMiddleware() proto.Middleware {
	if c.authMiddleware == nil {
		c.authMiddleware = auth_middleware.New(
			c.AuthService(),
			c.CookieManager(),
			c.Logger(),
		)
	}
	return c.authMiddleware
}

func (c *Container) RegistrationMiddleware() proto.Middleware {
	if c.registrationMiddleware == nil {
		c.registrationMiddleware = registration_middleware.New(
			c.AuthService(),
			c.CookieManager(),
		)
	}
	return c.registrationMiddleware
}

func (c *Container) RefreshTokenMiddleware() proto.Middleware {
	if c.refreshTokenMiddleware == nil {
		c.refreshTokenMiddleware = refresh_token_middleware.New(
			c.CookieManager(),
		)
	}

	return c.refreshTokenMiddleware
}

func (c *Container) AuthProviderCallbackGetHandler() *auth_provider_callback_post_handler.Handler {
	if c.authProviderCallbackGetHandler == nil {
		cfg := c.cfg.Auth
		c.authProviderCallbackGetHandler = auth_provider_callback_post_handler.New(
			c.AuthService(),
			c.OAuth(),
			c.CookieManager(),
			c.Logger(),
			cfg.LoggedInRedirectURL,
			cfg.RegistrationRedirectURL,
			cfg.ErrorRedirectURL,
		)
	}
	return c.authProviderCallbackGetHandler
}

func (c *Container) AuthProviderGetHandler() *auth_provider_get_handler.Handler {
	if c.authProviderGetHandler == nil {
		c.authProviderGetHandler = auth_provider_get_handler.New(
			c.OAuth(),
			c.Logger(),
		)
	}
	return c.authProviderGetHandler
}

func (c *Container) AuthLogoutPostHandler() *auth_logout_post_handler.Handler {
	if c.authLogoutPostHandler == nil {
		c.authLogoutPostHandler = auth_logout_post_handler.New(
			c.AuthService(),
			c.CookieManager(),
			c.OAuth(),
			c.Logger(),
		)
	}
	return c.authLogoutPostHandler
}

func (c *Container) AuthRefreshPostHandler() *auth_refresh_post_handler.Handler {
	if c.authRefreshPostHandler == nil {
		c.authRefreshPostHandler = auth_refresh_post_handler.New(
			c.AuthService(),
			c.CookieManager(),
			c.Logger(),
		)
	}
	return c.authRefreshPostHandler
}

func (c *Container) AuthRegisterPostHandler() *auth_register_post_handler.Handler {
	if c.authRegisterPostHandler == nil {
		c.authRegisterPostHandler = auth_register_post_handler.New(
			c.AuthService(),
			c.CookieManager(),
			c.Logger(),
		)
	}
	return c.authRegisterPostHandler
}

func (c *Container) UsersUsernameAvailabilityGetHandler() *users_username_availability_get_handler.Handler {
	if c.usersUsernameAvailabilityGetHandler == nil {
		c.usersUsernameAvailabilityGetHandler = users_username_availability_get_handler.New(
			c.AuthService(),
			c.Logger(),
		)
	}
	return c.usersUsernameAvailabilityGetHandler
}

func (c *Container) AuthPingGetHandler() *auth_ping_get_handler.Handler {
	if c.authPingGetHandler == nil {
		c.authPingGetHandler = auth_ping_get_handler.New()
	}
	return c.authPingGetHandler
}

func (c *Container) UsersMeStatisticsPostHandler() *users_me_statistics_post_handler.Handler {
	if c.usersMeStatisticsPostHandler == nil {
		c.usersMeStatisticsPostHandler = users_me_statistics_post_handler.New(
			c.StatisticsService(),
			c.Logger(),
		)
	}
	return c.usersMeStatisticsPostHandler
}

func (c *Container) UsersMeStatisticsGetHandler() *users_me_statistics_get_handler.Handler {
	if c.usersMeStatisticsGetHandler == nil {
		c.usersMeStatisticsGetHandler = users_me_statistics_get_handler.New(
			c.StatisticsService(),
			c.Logger(),
		)
	}
	return c.usersMeStatisticsGetHandler
}

func (c *Container) UsersMeStatisticsDeleteHandler() *users_me_statistics_delete_handler.Handler {
	if c.usersMeStatisticsDeleteHandler == nil {
		c.usersMeStatisticsDeleteHandler = users_me_statistics_delete_handler.New(
			c.Logger(),
			c.StatisticsService(),
		)
	}
	return c.usersMeStatisticsDeleteHandler
}

func (c *Container) UsersUsernameProfileGetHandler() *users_username_profile_get_handler.Handler {
	if c.usersUsernameProfileGetHandler == nil {
		c.usersUsernameProfileGetHandler = users_username_profile_get_handler.New(
			c.ProfileService(),
			c.Logger(),
		)
	}
	return c.usersUsernameProfileGetHandler
}

func (c *Container) UsersMeGetHandler() *users_me_get_handler.Handler {
	if c.usersMeGetHandler == nil {
		c.usersMeGetHandler = users_me_get_handler.New(
			c.UserService(),
			c.Logger(),
		)
	}
	return c.usersMeGetHandler
}

func (c *Container) AntifroadKeyGetHandler() *antifroad_key_get_handler.Handler {
	if c.antifroadKeyGetHandler == nil {
		c.antifroadKeyGetHandler = antifroad_key_get_handler.New(
			c.AntifroadService(),
			c.Logger(),
		)
	}
	return c.antifroadKeyGetHandler
}

func (c *Container) AntifroadRotateKeysPostHandler() *antifroad_rotate_keys_post_handler.Handler {
	if c.antifroadRotateKeysPostHandler == nil {
		c.antifroadRotateKeysPostHandler = antifroad_rotate_keys_post_handler.New(
			c.AntifroadService(),
			c.Logger(),
		)
	}
	return c.antifroadRotateKeysPostHandler
}

func (c *Container) UsersMeUsernamePatchHandler() *users_me_username_patch_handler.Handler {
	if c.usersMeUsernamePatchHandler == nil {
		c.usersMeUsernamePatchHandler = users_me_username_patch_handler.New(
			c.Logger(),
			c.UserService(),
		)
	}
	return c.usersMeUsernamePatchHandler
}
