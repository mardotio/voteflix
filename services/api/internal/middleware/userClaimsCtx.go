package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/internal/utils"
)

func UserClaimsCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			userClaims := utils.UserJwtClaims{
				Sub:  claims["sub"].(string),
				List: claims["list"].(string),
			}

			ctx := userClaims.WithValue(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(handler)
	}
}
