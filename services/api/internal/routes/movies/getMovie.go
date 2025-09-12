package movies

import (
	"context"
	"github.com/uptrace/bun"
	"net/http"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type owner struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatarUrl"`
}

type usersMap map[string]owner

type rating struct {
	Rating    int64                `json:"rating"`
	UserId    string               `json:"userId"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
}

type vote struct {
	IsApproval bool                 `json:"isApproval"`
	UserId     string               `json:"userId"`
	CreatedAt  utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt  *utils.JsonEpochTime `json:"updatedAt"`
}

type getMovieResponse struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	ListId    string               `json:"listId"`
	Status    string               `json:"status"`
	Votes     []vote               `json:"votes"`
	Ratings   []rating             `json:"ratings"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
	CreatorId string               `json:"creatorId"`
	Users     usersMap             `json:"users"`
}

func (body getMovieResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func getMovieRatings(db *bun.DB, ctx context.Context, movie models.Movie) ([]rating, error) {
	var ratings []models.Rating

	err := db.NewSelect().
		Model(&ratings).
		Where("movie_id = ?", movie.Id).
		Order("created_at asc").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	response := make([]rating, len(ratings))

	for i, r := range ratings {
		response[i] = rating{
			Rating:    r.Rating,
			UserId:    r.UserId,
			CreatedAt: utils.JsonEpochTime(r.CreatedAt),
			UpdatedAt: (*utils.JsonEpochTime)(r.UpdatedAt),
		}
	}

	return response, nil
}

func getMovieVotes(db *bun.DB, ctx context.Context, movie models.Movie) ([]vote, error) {
	var votes []models.Vote

	err := db.NewSelect().
		Model(&votes).
		Where("movie_id = ?", movie.Id).
		Order("created_at asc").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	response := make([]vote, len(votes))

	for i, v := range votes {
		response[i] = vote{
			IsApproval: v.IsApproval,
			UserId:     v.UserId,
			CreatedAt:  utils.JsonEpochTime(v.CreatedAt),
			UpdatedAt:  (*utils.JsonEpochTime)(v.UpdatedAt),
		}
	}

	return response, nil
}

type userWithNickname struct {
	models.User
	DiscordNickname *string
}

func getRelevantUsers(db *bun.DB, ctx context.Context, movie models.Movie, ratings []rating, votes []vote) (usersMap, error) {
	var users []userWithNickname
	presentUsers := make(map[string]bool)
	ids := make([]string, 0)

	for _, r := range ratings {
		presentUsers[r.UserId] = true
		ids = append(ids, r.UserId)
	}

	for _, v := range votes {
		if !presentUsers[v.UserId] {
			presentUsers[v.UserId] = true
			ids = append(ids, v.UserId)
		}
	}

	if !presentUsers[movie.CreatorId] {
		presentUsers[movie.CreatorId] = true
		ids = append(ids, movie.CreatorId)
	}

	userSubQuery := db.NewSelect().
		Model((*models.User)(nil)).
		Where("id IN (?)", bun.In(ids))

	err := db.NewSelect().
		TableExpr("(?) AS u", userSubQuery).
		Column("u.*", "lu.discord_nickname").
		Join("join list_users as lu").
		JoinOn("lu.user_id = u.id").
		JoinOn("lu.list_id = ?", movie.ListId).
		Scan(ctx, &users)

	if err != nil {
		return nil, err
	}

	userMap := make(usersMap)

	for _, u := range users {
		name := u.DiscordUsername

		if u.DiscordNickname != nil {
			name = *u.DiscordNickname
		}

		userMap[u.Id] = owner{
			Name:      name,
			AvatarUrl: utils.GetAvatarUrl(u.DiscordId, u.DiscordAvatarId),
		}
	}

	return userMap, nil
}

func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)

	targetMovie := middleware.GetMovieFromCtx(ctx)
	ratings, ratingsErr := getMovieRatings(db, ctx, targetMovie)

	if ratingsErr != nil {
		jsonSender.InternalServerError(ratingsErr)
		return
	}

	votes, votesErr := getMovieVotes(db, ctx, targetMovie)

	if votesErr != nil {
		jsonSender.InternalServerError(votesErr)
		return
	}

	users, usersErr := getRelevantUsers(db, ctx, targetMovie, ratings, votes)

	if usersErr != nil {
		jsonSender.InternalServerError(usersErr)
		return
	}

	jsonSender.Ok(getMovieResponse{
		Id:        targetMovie.Id,
		Name:      targetMovie.Name,
		ListId:    targetMovie.ListId,
		Status:    targetMovie.Status,
		Ratings:   ratings,
		Votes:     votes,
		CreatorId: targetMovie.CreatorId,
		CreatedAt: utils.JsonEpochTime(targetMovie.CreatedAt),
		UpdatedAt: (*utils.JsonEpochTime)(targetMovie.UpdatedAt),
		Users:     users,
	})
}
