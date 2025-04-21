package util

import (
	"context"

	"github.com/google/uuid"
)

const (
	UserIDCTX    = "user_id"
	UserEmailCTX = "email"
)

func UserIDFromCTX(ctx context.Context) (userID uuid.UUID) {
	userID, ok := ctx.Value(UserIDCTX).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil
	}

	return userID
}

func UserEmailFromCTX(ctx context.Context) (email string) {
	user := ctx.Value(UserEmailCTX)
	email, _ = user.(string)
	return
}
