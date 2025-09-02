package movies

import "voteflix/api/internal/app"

type Handler struct{ app *app.App }

func NewMoviesHAndler(app *app.App) *Handler { return &Handler{app: app} }
