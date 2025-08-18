package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"voteflix/api/internal/app"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/routes/auth"
	"voteflix/api/internal/routes/bot"
	"voteflix/api/internal/routes/users"
)

func Router(app *app.App) {
	r := app.Router()

	authHandler := auth.NewAuthHandler(app)
	usersHandler := users.NewUsersHandler(app)
	botHandler := bot.NewBotHandler(app)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.JwtAuth()))
		r.Use(middleware.Authenticator())

		r.Route("/users", func(r chi.Router) {
			r.Get("/whoami", usersHandler.WhoAmI)
		})
	})

	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/token", authHandler.CreateToken)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.BotAuthenticator(app.Config()))

		r.Route("/bot", func(r chi.Router) {
			r.Post("/list", botHandler.CreateList)
		})
	})
}
