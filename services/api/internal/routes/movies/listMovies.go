package movies

import (
	"context"
	"github.com/uptrace/bun"
	"net/http"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type listMovieCreator struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatarUrl"`
}

type listMovieDetails struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	Creator   listMovieCreator     `json:"creator"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
}

type listMoviesResponse struct {
	Data []listMovieDetails `json:"data"`
}

func (body listMoviesResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

type movieWithUser struct {
	models.Movie
	DiscordNickname *string
	DiscordAvatarId *string
	DiscordUsername string
}

func listMovies(ctx context.Context, db *bun.DB, claims utils.UserJwtClaims) ([]movieWithUser, error) {
	var movies []movieWithUser

	moviesSubQuery := db.NewSelect().
		Model((*models.Movie)(nil)).
		Where("list_id = ?", claims.Scope).
		Limit(10)

	err := db.NewSelect().
		TableExpr("(?) AS movie", moviesSubQuery).
		Column("movie.*", "c.discord_nickname", "u.discord_avatar_id", "u.discord_username").
		Join("join list_users as c").
		JoinOn("c.user_id = movie.creator_id").
		JoinOn("c.list_id = movie.list_id").
		Join("join users as u").
		JoinOn("u.id = movie.creator_id").
		Scan(ctx, &movies)

	return movies, err
}

func toListMoviesResponse(movies []movieWithUser) listMoviesResponse {
	response := make([]listMovieDetails, len(movies))

	for i, m := range movies {
		creator := listMovieCreator{
			Name:      m.DiscordUsername,
			AvatarUrl: utils.GetAvatarUrl(models.User{DiscordAvatarId: m.DiscordAvatarId}),
		}

		if m.DiscordNickname != nil {
			creator.Name = *m.DiscordNickname
		}

		response[i] = listMovieDetails{
			Id:        m.Id,
			Name:      m.Name,
			Status:    m.Status,
			Creator:   creator,
			CreatedAt: utils.JsonEpochTime(m.CreatedAt),
			UpdatedAt: (*utils.JsonEpochTime)(m.UpdatedAt),
		}
	}

	return listMoviesResponse{Data: response}
}

func (h *Handler) ListMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := middleware.GetUserClaimsFromCtx(ctx)

	movies, moviesErr := listMovies(ctx, db, userClaims)

	if moviesErr != nil {
		jsonSender.InternalServerError(moviesErr)
		return
	}

	jsonSender.Ok(toListMoviesResponse(movies))
}
