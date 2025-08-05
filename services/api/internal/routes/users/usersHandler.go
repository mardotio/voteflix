package users

import "voteflix/api/internal/app"

type Handler struct {
	app *app.App
}

func NewUsersHandler(app *app.App) *Handler {
	return &Handler{app: app}
}
