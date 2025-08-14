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
			errorSender := utils.NewJsonSender(w, r)

			if err != nil {
				errorSender.Unauthorized(err)
				return
			}

			if token == nil {
				errorSender.Unauthorized(fmt.Errorf("not token found"))
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
