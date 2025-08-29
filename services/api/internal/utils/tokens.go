package utils

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"time"
)

type userCtxKey string

const userClaimsCtx userCtxKey = "userClaims"

type UserJwtClaims struct {
	Sub   string
	Scope string
}

func (claims UserJwtClaims) ToClaimsMap(duration time.Duration) map[string]interface{} {
	claimsMap := map[string]interface{}{
		"sub":   claims.Sub,
		"scope": claims.Scope,
	}
	jwtauth.SetIssuedNow(claimsMap)
	jwtauth.SetExpiryIn(claimsMap, duration)

	return claimsMap
}

func (claims UserJwtClaims) WithValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, userClaimsCtx, claims)
}

func GetUserClaimsFromCtx(ctx context.Context) UserJwtClaims {
	return ctx.Value(userClaimsCtx).(UserJwtClaims)
}

type botCtxKey string

const botClaimsCtx botCtxKey = "botClaims"

type BotJwtClaims struct {
	Sub      string
	Server   string
	Username string
	Avatar   *string
	Nickname *string
}

func (claims BotJwtClaims) ToClaimsMap(duration time.Duration) map[string]interface{} {
	claimsMap := map[string]interface{}{
		"sub":      claims.Sub,
		"server":   claims.Server,
		"username": claims.Username,
	}

	if claims.Avatar != nil {
		claimsMap["avatar"] = *claims.Avatar
	}

	if claims.Nickname != nil {
		claimsMap["nickname"] = *claims.Nickname
	}

	jwtauth.SetIssuedNow(claimsMap)
	jwtauth.SetExpiryIn(claimsMap, duration)

	return claimsMap
}

func (claims BotJwtClaims) WithValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, botClaimsCtx, claims)
}

func GetBotClaimsFromCtx(ctx context.Context) BotJwtClaims {
	return ctx.Value(botClaimsCtx).(BotJwtClaims)
}
