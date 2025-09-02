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

type moviePathParams struct {
	MovieId string `validate:"required,uuid4"`
}

type movieCtxKey string

const movieCtx movieCtxKey = "movie"

func withValue(ctx context.Context, movie models.Movie) context.Context {
	return context.WithValue(ctx, movieCtx, movie)
}

func GetMovieFromCtx(ctx context.Context) models.Movie {
	return ctx.Value(movieCtx).(models.Movie)
}

func MovieCtx(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			v := app.Validate()
			db := app.Db()
			jsonSender := utils.NewJsonSender(w, r)

			userClaims := utils.GetUserClaimsFromCtx(ctx)
			params := moviePathParams{
				MovieId: chi.URLParam(r, "movieId"),
			}

			if err := v.Struct(params); err != nil {
				jsonSender.BadRequest(err)
				return
			}

			movie := models.Movie{}
			movieErr := db.NewSelect().
				Model(&movie).
				Where("id = ?", params.MovieId).
				Where("list_id = ?", userClaims.Scope).
				Scan(ctx)

			if errors.Is(movieErr, sql.ErrNoRows) {
				jsonSender.NotFound(fmt.Errorf("could not find movie with id %s", params.MovieId))
				return
			}

			if movieErr != nil {
				jsonSender.InternalServerError(movieErr)
				return
			}

			next.ServeHTTP(w, r.WithContext(withValue(ctx, movie)))
		}

		return http.HandlerFunc(handler)
	}
}
