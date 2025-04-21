//nolint:nolintlint,staticcheck
package tracking

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type ContextKey string

const (
	KeyContextID ContextKey = "context_id"
)

func GetTrackIDFromContext(ctx context.Context) string {
	return cast.ToString(ctx.Value(KeyContextID))
}

func CloneTrackeIDToCtx(fromCtx context.Context, toCtx context.Context) context.Context {
	traceID := GetTrackIDFromContext(fromCtx)
	return context.WithValue(toCtx, KeyContextID, traceID)
}

func InitContextWithTrackID() context.Context {
	var trackID = GenTrackID()
	return context.WithValue(context.Background(), KeyContextID, &trackID) // nolint: staticcheck
}

func GenTrackID() string {
	return fmt.Sprintf("wheel-%s", uuid.New().String())
}
