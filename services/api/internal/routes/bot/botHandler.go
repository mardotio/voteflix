package bot

import "voteflix/api/internal/app"

type Handler struct{ app *app.App }

func NewBotHandler(app *app.App) *Handler { return &Handler{app: app} }
