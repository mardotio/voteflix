package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"log"
	"math/rand"
	"net/http"
	"time"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type pickMovieCreator struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatarUrl"`
}

type pickMovieResponse struct {
	Id        string              `json:"id"`
	Name      string              `json:"name"`
	Creator   pickMovieCreator    `json:"creator"`
	CreatedAt utils.JsonEpochTime `json:"createdAt"`
}

func (body pickMovieResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

type orderer struct {
	comparator bun.Safe
	direction  bun.Safe
}

func getDirection() (int, orderer, orderer) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	left := orderer{comparator: "<=", direction: "desc"}
	right := orderer{comparator: ">", direction: "asc"}

	seed := r.Intn(2147483648)

	if r.Intn(2) == 0 {
		seed = seed * -1
	}

	if r.Intn(2) == 0 {
		left.direction = "asc"
		right.direction = "desc"
	}

	if r.Intn(2) == 0 {
		return seed, left, right
	}

	return seed, right, left
}

func findMovie(db *bun.DB, ctx context.Context, list models.List) (models.Movie, error) {
	var movie models.Movie

	seed, firstDir, secondDir := getDirection()

	firstDirErr := db.NewSelect().
		Model(&movie).
		Relation("Creator").
		Relation("User").
		Where("status = ?", "approved").
		Where("movie.list_id = ?", list.Id).
		Where("seed ? ?", firstDir.comparator, seed).
		OrderExpr("? ?", bun.Ident("seed"), firstDir.direction).
		Limit(1).
		Scan(ctx)

	if firstDirErr != nil && !errors.Is(firstDirErr, sql.ErrNoRows) {
		return movie, firstDirErr
	}

	if firstDirErr != nil {
		secondDirErr := db.NewSelect().
			Model(&movie).
			Relation("Creator").
			Relation("User").
			Where("status = ?", "approved").
			Where("movie.list_id = ?", list.Id).
			Where("seed ? ?", secondDir.comparator, seed).
			OrderExpr("? ?", bun.Ident("seed"), secondDir.direction).
			Limit(1).
			Scan(ctx)

		if secondDirErr != nil {
			return movie, secondDirErr
		}
	}

	_, seedUpdateErr := db.NewUpdate().
		Model(&movie).
		Where("id = ?", movie.Id).
		Set("seed = gen_random_movie_seed()").
		Exec(ctx)

	return movie, seedUpdateErr
}

func getHasPendingMovie(db *bun.DB, ctx context.Context, list models.List) (bool, error) {
	return db.NewSelect().
		Model((*models.Movie)(nil)).
		Where("status = 'pending'").
		Where("list_id = ?", list.Id).
		Exists(ctx)
}

func (h *Handler) PickMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)
	server := middleware.GetServerFromCtx(ctx)

	hasPendingMovie, hasPendingMovieErr := getHasPendingMovie(db, ctx, server)

	if hasPendingMovieErr != nil {
		jsonSender.InternalServerError(hasPendingMovieErr)
		return
	}

	if hasPendingMovie {
		jsonSender.UnprocessableEntity(fmt.Errorf("cannot pick movie if there are movies pending in list"))
		return
	}

	movie, err := findMovie(db, ctx, server)

	if errors.Is(err, sql.ErrNoRows) {
		jsonSender.NotFound(err)
		return
	}

	if err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	log.Println("seed", movie.Seed, movie.Name)

	name := movie.User.DiscordUsername

	if movie.Creator.DiscordNickname != nil {
		name = *movie.Creator.DiscordNickname
	}

	jsonSender.Ok(pickMovieResponse{
		Id:   movie.Id,
		Name: movie.Name,
		Creator: pickMovieCreator{
			Name:      name,
			AvatarUrl: utils.GetAvatarUrl(movie.User.DiscordId, movie.User.DiscordAvatarId, false),
		},
		CreatedAt: utils.JsonEpochTime(movie.CreatedAt),
	})
}
