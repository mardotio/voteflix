package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"voteflix/api/internal/app"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type serverCtxKey string

const serverCtx serverCtxKey = "server"

func withServerCtx(ctx context.Context, list models.List) context.Context {
	return context.WithValue(ctx, serverCtx, list)
}

func GetServerFromCtx(ctx context.Context) models.List {
	return ctx.Value(serverCtx).(models.List)
}

func ServerCtx(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			db := app.Db()
			jsonSender := utils.NewJsonSender(w, r)

			serverId := chi.URLParam(r, "serverId")

			if serverId == "" {
				jsonSender.BadRequest(errors.New("serverId is required"))
				return
			}

			list := models.List{}
			listErr := db.NewSelect().
				Model(&list).
				Where("discord_server_id = ?", serverId).
				Scan(ctx)

			if errors.Is(listErr, sql.ErrNoRows) {
				jsonSender.NotFound(fmt.Errorf("could not find list with id %s", serverId))
				return
			}

			if listErr != nil {
				jsonSender.InternalServerError(listErr)
				return
			}

			next.ServeHTTP(w, r.WithContext(withServerCtx(ctx, list)))
		}

		return http.HandlerFunc(handler)
	}
}
