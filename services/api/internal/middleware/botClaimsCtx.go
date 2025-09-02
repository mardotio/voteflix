package middleware

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/internal/utils"
)

type botCtxKey string

const botClaimsCtx botCtxKey = "botClaims"

func withBotClaimsCtx(ctx context.Context, claims utils.BotJwtClaims) context.Context {
	return context.WithValue(ctx, botClaimsCtx, claims)
}

func GetBotClaimsFromCtx(ctx context.Context) utils.BotJwtClaims {
	return ctx.Value(botClaimsCtx).(utils.BotJwtClaims)
}

func BotClaimsCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			botClaims := utils.BotJwtClaims{
				Sub:      claims["sub"].(string),
				Server:   claims["server"].(string),
				Username: claims["username"].(string),
			}

			if avatar, ok := claims["avatar"].(string); ok && avatar != "" {
				botClaims.Avatar = &avatar
			}

			if nickname, ok := claims["nickname"].(string); ok && nickname != "" {
				botClaims.Nickname = &nickname
			}

			next.ServeHTTP(w, r.WithContext(withBotClaimsCtx(r.Context(), botClaims)))
		}

		return http.HandlerFunc(handler)
	}
}
