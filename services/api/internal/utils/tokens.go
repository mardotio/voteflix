package utils

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
	"voteflix/api/internal/app"
)

type AppToken interface {
	toClaimsMap(duration time.Duration) map[string]interface{}
	jwtAuth(app *app.App) *jwtauth.JWTAuth
}

func GetAppToken(app *app.App, token AppToken, duration time.Duration) (jwt.Token, string, error) {
	return token.jwtAuth(app).Encode(token.toClaimsMap(duration))
}

type UserJwtClaims struct {
	Sub   string
	Scope string
}

func (claims UserJwtClaims) jwtAuth(app *app.App) *jwtauth.JWTAuth {
	return app.JwtAuth()
}

func (claims UserJwtClaims) toClaimsMap(duration time.Duration) map[string]interface{} {
	claimsMap := map[string]interface{}{
		"sub":   claims.Sub,
		"scope": claims.Scope,
	}
	jwtauth.SetIssuedNow(claimsMap)
	jwtauth.SetExpiryIn(claimsMap, duration)

	return claimsMap
}

type BotJwtClaims struct {
	Sub      string
	Server   string
	Username string
	Avatar   *string
	Nickname *string
}

func (claims BotJwtClaims) jwtAuth(app *app.App) *jwtauth.JWTAuth {
	return app.BotJwtAuth()
}

func (claims BotJwtClaims) toClaimsMap(duration time.Duration) map[string]interface{} {
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
