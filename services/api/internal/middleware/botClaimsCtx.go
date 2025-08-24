package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/internal/utils"
)

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

			ctx := botClaims.WithValue(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(handler)
	}
}
