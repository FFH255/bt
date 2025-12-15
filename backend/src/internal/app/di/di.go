package di

import (
	"github.com/robfig/cron"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api/cookie"
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
	"github.com/ruslanonly/blindtyping/src/internal/app/config"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/antifroad_key_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/blocked_token_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/language_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/pb_cache"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/profiles_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/session_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
	"github.com/ruslanonly/blindtyping/src/internal/scheduler/handlers/antifroad_rotate_keys_handler"
	"github.com/ruslanonly/blindtyping/src/internal/scheduler/handlers/expired_sessions_handler"
	"github.com/ruslanonly/blindtyping/src/internal/services/antifroad_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/auth_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/pb_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/profile_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/session_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/statistics_service"
	"github.com/ruslanonly/blindtyping/src/internal/services/user_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/oauth"
	"github.com/ruslanonly/blindtyping/src/internal/shared/postgres"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
	"github.com/ruslanonly/blindtyping/src/internal/shared/redis"
	"github.com/ruslanonly/blindtyping/src/internal/shared/uuid_generator"
)

type Container struct {
	logger        internal.Logger
	cfg           *config.Config
	postgres      *postgres.Database
	redis         *redis.Client
	uuidGenerator *uuid_generator.Generator
	oauth         *oauth.OAuth
	cookieManager *cookie.Manager
	// Repositories
	userRepository         *user_repository.Repository
	sessionRepository      *session_repository.Repository
	blockedTokenRepository *blocked_token_repository.Repository
	statisticsRepository   *statistics_repository.Repository
	profileRepository      *profiles_repository.Repository
	pbCache                *pb_cache.Cache
	languageRepository     *language_repository.Repository
	antifroadKeyRepository *antifroad_key_repository.Repository
	// Services
	sessionService        *session_service.Service
	authService           *auth_service.Service
	userService           *user_service.Service
	statisticsService     *statistics_service.Service
	pbService             *pb_service.Service
	profileService        *profile_service.Service
	antifroadKeyGenerator *antifroad_service.KeyGenerator
	antifroadService      *antifroad_service.Service
	// Handlers
	authProviderCallbackGetHandler      *auth_provider_callback_post_handler.Handler
	authProviderGetHandler              *auth_provider_get_handler.Handler
	authLogoutPostHandler               *auth_logout_post_handler.Handler
	authRefreshPostHandler              *auth_refresh_post_handler.Handler
	authPingGetHandler                  *auth_ping_get_handler.Handler
	authRegisterPostHandler             *auth_register_post_handler.Handler
	usersUsernameAvailabilityGetHandler *users_username_availability_get_handler.Handler
	usersMeStatisticsPostHandler        *users_me_statistics_post_handler.Handler
	usersMeStatisticsGetHandler         *users_me_statistics_get_handler.Handler
	usersMeStatisticsDeleteHandler      *users_me_statistics_delete_handler.Handler
	usersUsernameProfileGetHandler      *users_username_profile_get_handler.Handler
	usersMeGetHandler                   *users_me_get_handler.Handler
	usersMeUsernamePatchHandler         *users_me_username_patch_handler.Handler
	antifroadKeyGetHandler              *antifroad_key_get_handler.Handler
	antifroadRotateKeysPostHandler      *antifroad_rotate_keys_post_handler.Handler
	//Middleware
	authMiddleware         proto.Middleware
	registrationMiddleware proto.Middleware
	refreshTokenMiddleware proto.Middleware
	// Server
	router *proto.Router
	server *proto.Server
	// Scheduler
	scheduler                  *cron.Cron
	expiredSessionsHandler     *expired_sessions_handler.Handler
	antifroadRotateKeysHandler *antifroad_rotate_keys_handler.Handler
}

func (c *Container) Server() *proto.Server {
	if c.server == nil {
		cfg := c.cfg.Server

		c.server = proto.NewServer(
			c.Router(),
			cfg.Address,
			proto.MustUnmarshalDuration(cfg.ShutdownTimeout),
			proto.MustUnmarshalDuration(cfg.ReadTimeout),
			proto.MustUnmarshalDuration(cfg.WriteTimeout),
			proto.MustUnmarshalDuration(cfg.IdleTimeout),
		)
	}
	return c.server
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}
