package internal

import "context"

type Logger interface {
	Debug(ctx context.Context, args ...any)
	Info(ctx context.Context, args ...any)
	Warning(ctx context.Context, args ...any)
	Error(ctx context.Context, args ...any)
	WithField(ctx context.Context, k string, v any) context.Context
	WithFields(ctx context.Context, fields map[string]any) context.Context
	WithError(ctx context.Context, err error) context.Context
	WithRequestID(ctx context.Context, requestID string) context.Context
	WithUserID(ctx context.Context, userID int64) context.Context
	WithHandlerName(ctx context.Context, name string) context.Context
	WithStatusCode(ctx context.Context, statusCode int) context.Context
	WithMsg(ctx context.Context, msg string) context.Context
}
