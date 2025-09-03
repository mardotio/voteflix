package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"voteflix/api/internal/app"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/routes/auth"
	"voteflix/api/internal/routes/bot"
	"voteflix/api/internal/routes/movies"
	"voteflix/api/internal/routes/users"
)

func Router(app *app.App) {
	r := app.Router()

	authHandler := auth.NewAuthHandler(app)
	usersHandler := users.NewUsersHandler(app)
	botHandler := bot.NewBotHandler(app)
	moviesHandler := movies.NewMoviesHAndler(app)

	//Auth related routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.BotJwtAuth()))
		r.Use(middleware.Authenticator())
		r.Use(middleware.BotClaimsCtx())

		r.Route("/auth", func(r chi.Router) {
			r.Post("/token", authHandler.CreateToken)
		})
	})

	//Routes accessed by discord bot
	r.Group(func(r chi.Router) {
		r.Use(middleware.BotAuthenticator(app.Config()))

		r.Route("/bot", func(r chi.Router) {
			if !app.Config().IsProduction() {
				r.Post("/token", botHandler.CreateToken)
			}
			r.Post("/list", botHandler.CreateList)
		})
	})

	//Rotes for main application
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.JwtAuth()))
		r.Use(middleware.Authenticator())
		r.Use(middleware.UserClaimsCtx())

		r.Route("/users", func(r chi.Router) {
			r.Get("/whoami", usersHandler.WhoAmI)
		})

		r.Route("/movies", func(r chi.Router) {
			r.Post("/", moviesHandler.CreateMovie)

			r.Route("/{movieId}", func(r chi.Router) {
				r.Use(middleware.MovieCtx(app))

				r.Patch("/", moviesHandler.UpdateMovie)

				r.Route("/votes", func(r chi.Router) {
					r.Put("/", moviesHandler.AddMovieVote)
				})

				r.Route("/ratings", func(r chi.Router) {
					r.Put("/", moviesHandler.AddMovieRating)
				})
			})
		})
	})
}
