package movies

import (
	"database/sql"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type addMovieRatingRequest struct {
	Rating *int64 `json:"rating" validate:"required,min=0,max=10"`

	movie     models.Movie
	validator *validator.Validate
}

func (body addMovieRatingRequest) Bind(r *http.Request) error {
	if err := body.validator.Struct(body); err != nil {
		return err
	}

	if body.movie.Status != "watched" {
		return errors.New("cannot add rating to movie that has not been watched")
	}

	return nil
}

type addMovieRatingResponse struct {
	MovieId   string               `json:"movieId"`
	Rating    int64                `json:"rating"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
}

func (body addMovieRatingResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) AddMovieRating(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	v := h.app.Validate()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := middleware.GetUserClaimsFromCtx(ctx)
	targetMovie := middleware.GetMovieFromCtx(ctx)
	body := addMovieRatingRequest{validator: v, movie: targetMovie}

	if err := render.Bind(r, &body); err != nil {
		jsonSender.BadRequest(err)
		return
	}

	now := time.Now()
	rating := models.Rating{
		MovieId: targetMovie.Id,
		UserId:  userClaims.Sub,
		Rating:  *body.Rating,
		Timestamps: models.Timestamps{
			CreatedAt: now,
		},
	}

	err := db.NewInsert().
		Model(&rating).
		On("conflict (user_id, movie_id) do update").
		Set("rating = ?", rating.Rating).
		Set("updated_at = ?", now).
		Where("rating.rating != ?", rating.Rating).
		Returning("*").
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		jsonSender.Conflict(errors.New("rating with same value already exists for this movie"))
		return
	}

	jsonSender.Ok(addMovieRatingResponse{
		MovieId:   targetMovie.Id,
		Rating:    rating.Rating,
		CreatedAt: utils.JsonEpochTime(rating.CreatedAt),
		UpdatedAt: (*utils.JsonEpochTime)(rating.UpdatedAt),
	})
}
