package auth

import (
	"context"

	"github.com/gofrs/uuid"
)

type contextKey string

const userIDKey = contextKey("UserID")

// WithUserID(ctx context.Context, userID uuid.UUID) context.Context
// добавляет в контекст ID пользователя текущего
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFrom(ctx context.Context) (uuid.UUID, bool)
// достаёт из контекста ID пользователя текущего
func UserIDFrom(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	return v, ok
}
