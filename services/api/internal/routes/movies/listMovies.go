package movies

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"net/http"
	"net/url"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type listMoviesQueryParams struct {
	Status string `validate:"omitempty,oneof=pending approved watched rejected"`
	Query  string `validate:"omitempty,max=100"`
}

func getListMoviesQueryParams(values url.Values, validate *validator.Validate) (listMoviesQueryParams, error) {
	params := listMoviesQueryParams{
		Status: values.Get("status"),
		Query:  values.Get("query"),
	}

	err := validate.Struct(&params)

	return params, err
}

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
	WatchedAt *utils.JsonEpochTime `json:"watchedAt"`
}

func (r listMovieDetails) CursorId() *string { return &r.Id }

type movieWithUser struct {
	models.Movie
	DiscordNickname *string
	DiscordAvatarId *string
	DiscordUsername string
	DiscordId       string
}

func listMovies(
	ctx context.Context,
	db *bun.DB,
	claims utils.UserJwtClaims,
	cursor utils.Cursor[listMovieDetails],
	queryParams listMoviesQueryParams,
) ([]movieWithUser, error) {
	var movies []movieWithUser

	cursorOrder := cursor.Order()

	moviesSubQuery := db.NewSelect().
		Model((*models.Movie)(nil)).
		Where("list_id = ?", claims.Scope).
		OrderExpr("? ?", bun.Ident("created_at"), cursorOrder).
		OrderExpr("? ?", bun.Ident("id"), cursorOrder).
		Limit(cursor.FetchLimit())

	if queryParams.Status != "" {
		moviesSubQuery.Where("status = ?", queryParams.Status)
	}

	if queryParams.Query != "" {
		moviesSubQuery.Where("name ilike ?", fmt.Sprintf("%%%s%%", queryParams.Query))
	}

	if !cursor.IsStart() {
		subQuery := db.NewSelect().
			Model((*models.Movie)(nil)).
			Column("created_at").
			Where("list_id = ?", claims.Scope).
			Where("id = ?", *cursor.Marker())

		moviesSubQuery.
			Where("(created_at, id) ? ((?), ?)", cursor.Comparator(), subQuery, *cursor.Marker())
	}

	err := db.NewSelect().
		TableExpr("(?) AS movie", moviesSubQuery).
		Column("movie.*", "c.discord_nickname", "u.discord_avatar_id", "u.discord_username", "u.discord_id").
		Join("join list_users as c").
		JoinOn("c.user_id = movie.creator_id").
		JoinOn("c.list_id = ?", claims.Scope).
		Join("join users as u").
		JoinOn("u.id = movie.creator_id").
		Scan(ctx, &movies)

	return movies, err
}

func toMovieDetails(movies []movieWithUser) []listMovieDetails {
	data := make([]listMovieDetails, len(movies))

	for i, m := range movies {
		creator := listMovieCreator{
			Name:      m.DiscordUsername,
			AvatarUrl: utils.GetAvatarUrl(m.DiscordId, m.DiscordAvatarId, false),
		}

		if m.DiscordNickname != nil {
			creator.Name = *m.DiscordNickname
		}

		data[i] = listMovieDetails{
			Id:        m.Id,
			Name:      m.Name,
			Status:    m.Status,
			Creator:   creator,
			CreatedAt: utils.JsonEpochTime(m.CreatedAt),
			UpdatedAt: (*utils.JsonEpochTime)(m.UpdatedAt),
			WatchedAt: (*utils.JsonEpochTime)(m.WatchedAt),
		}
	}

	return data
}

func (h *Handler) ListMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	v := h.app.Validate()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := middleware.GetUserClaimsFromCtx(ctx)
	queryParams, queryParamsErr := getListMoviesQueryParams(r.URL.Query(), v)

	if queryParamsErr != nil {
		jsonSender.BadRequest(queryParamsErr)
		return
	}

	cursor, cursorErr := utils.NewCursorFromMap[listMovieDetails](r.URL.Query(), v)

	if cursorErr != nil {
		jsonSender.BadRequest(cursorErr)
		return
	}

	movies, moviesErr := listMovies(ctx, db, userClaims, cursor, queryParams)

	if moviesErr != nil {
		jsonSender.InternalServerError(moviesErr)
		return
	}

	jsonSender.Ok(cursor.WithData(toMovieDetails(movies)).ToResponse())
}
