package middleware

import (
	"context"
	"errors"
	"net/http"
)

type HTTPMiddleware func(http.Handler) http.Handler

func GetUserCtx(ctx context.Context) (*UserCtx, error) {
	val := ctx.Value(User)
	userCtx, ok := val.(UserCtx)

	if !ok {
		return nil, errors.New("cannot get user context")
	}

	return &userCtx, nil
}

func Must(ctx *UserCtx, err error) *UserCtx {
	if err != nil {
		panic(err)
	}

	return ctx
}
