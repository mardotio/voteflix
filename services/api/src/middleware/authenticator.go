package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/src/utils"
)

type authErrorResponse struct {
	Message string `json:"message"`
}

func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				utils.JsonError(w, authErrorResponse{Message: err.Error()}, http.StatusUnauthorized)
				return
			}

			if token == nil {
				utils.JsonError(w, authErrorResponse{Message: "no token found"}, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
