package middleware

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/internal/utils"
)

type userCtxKey string

const userClaimsCtx userCtxKey = "userClaims"

func withUserClaimsCtx(ctx context.Context, claims utils.UserJwtClaims) context.Context {
	return context.WithValue(ctx, userClaimsCtx, claims)
}

func GetUserClaimsFromCtx(ctx context.Context) utils.UserJwtClaims {
	return ctx.Value(userClaimsCtx).(utils.UserJwtClaims)
}

func UserClaimsCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			userClaims := utils.UserJwtClaims{
				Sub:   claims["sub"].(string),
				Scope: claims["scope"].(string),
			}

			next.ServeHTTP(w, r.WithContext(withUserClaimsCtx(r.Context(), userClaims)))
		}

		return http.HandlerFunc(handler)
	}
}
