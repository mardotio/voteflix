package auth

import "voteflix/api/internal/app"

type Handler struct{ app *app.App }

func NewAuthHandler(app *app.App) *Handler {
	return &Handler{app: app}
}
