package api

import "context"

type ctxKey string

const ctxChannelID ctxKey = "channelID"

func withChannelID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxChannelID, id)
}

func channelIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(ctxChannelID).(string)
	return v
}
