package movies

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type updateMovieRequest struct {
	Name   *string `json:"name" validate:"omitempty,max=255,min=1,required_if=Status nil"`
	Status *string `json:"status" validate:"omitempty,oneof=watched approved,required_if=Name nil"`

	movie     models.Movie
	validator *validator.Validate
}

var noUpdateError = errors.New("No fields to update")

func (body *updateMovieRequest) Bind(r *http.Request) error {
	if err := body.validator.Struct(body); err != nil {
		return err
	}

	movie := body.movie

	if body.Name != nil && *body.Name == movie.Name {
		body.Name = nil
	}

	if body.Status != nil && *body.Status == movie.Status {
		body.Status = nil
	}

	if body.Name == nil && body.Status == nil {
		return noUpdateError
	}

	if utils.IsEqual(body.Name, &movie.Name) && utils.IsEqual(body.Status, &movie.Status) {
		return errors.New("")
	}

	if movie.Status == "rejected" {
		return errors.New("cannot update movie in \"rejected\" status")
	}

	if movie.Status == "pending" {
		if body.Status != nil {
			return errors.New("cannot update status of \"pending\" movie")
		}

		if movie.ApproveCount > 0 || movie.RejectCount > 0 {
			return errors.New("cannot update name or movie that has already been reviewed")
		}
		return nil
	}

	if body.Name != nil {
		return fmt.Errorf("cannot update name of \"%s\" movie", movie.Status)
	}

	return nil
}

type updateMovieResponse struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	Status    string               `json:"status"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
}

func (body updateMovieResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	v := h.app.Validate()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)

	targetMovie := middleware.GetMovieFromCtx(ctx)
	body := updateMovieRequest{validator: v, movie: targetMovie}
	validationError := render.Bind(r, &body)

	response := updateMovieResponse{
		Id:        targetMovie.Id,
		Name:      targetMovie.Name,
		Status:    targetMovie.Status,
		CreatedAt: utils.JsonEpochTime(targetMovie.CreatedAt),
		UpdatedAt: (*utils.JsonEpochTime)(targetMovie.UpdatedAt),
	}

	if errors.Is(validationError, noUpdateError) {
		jsonSender.Ok(response)
		return
	}

	if validationError != nil {
		jsonSender.BadRequest(validationError)
		return
	}

	now := time.Now()
	response.UpdatedAt = (*utils.JsonEpochTime)(&now)

	updateQuery := db.NewUpdate().
		Model((*models.Movie)(nil)).
		Where("id = ?", targetMovie.Id).
		Where("status = ?", targetMovie.Status).
		Set("updated_at = ?", now)

	if body.Name != nil {
		response.Name = *body.Name
		updateQuery.Set("name = ?", *body.Name)
	}

	if body.Status != nil {
		response.Status = *body.Status
		updateQuery.Set("status = ?", *body.Status)
	}

	_, updateErr := updateQuery.Exec(ctx)

	if updateErr != nil {
		jsonSender.InternalServerError(updateErr)
		return
	}

	jsonSender.Ok(response)
}
