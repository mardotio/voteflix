package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"net/http"
	customMiddleware "voteflix/api/src/middleware"
	"voteflix/api/src/routes/auth"
	"voteflix/api/src/routes/users"
	"voteflix/api/src/utils"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(utils.GetTokenAuth()))
		r.Use(customMiddleware.Authenticator())

		r.Route("/users", func(r chi.Router) {
			r.Get("/whoami", users.WhoAmI)
		})
	})

	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/token", auth.CreateToken)
		})
	})

	return r
}
