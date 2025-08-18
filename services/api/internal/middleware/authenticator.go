package middleware

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"voteflix/api/internal/utils"
)

func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			jsonSender := utils.NewJsonSender(w, r)

			if err != nil {
				jsonSender.Unauthorized(err)
				return
			}

			if token == nil {
				jsonSender.Unauthorized(fmt.Errorf("no token found"))
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
