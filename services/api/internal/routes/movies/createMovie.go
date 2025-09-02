package movies

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type createMovieResponse struct {
	Id        string              `json:"id"`
	ListId    string              `json:"listId"`
	Name      string              `json:"name"`
	Status    string              `json:"status"`
	CreatedAt utils.JsonEpochTime `json:"createdAt"`
}

func (body createMovieResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

type createMovieRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	validator *validator.Validate
}

func (body createMovieRequest) Bind(r *http.Request) error {
	return body.validator.Struct(body)
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	v := h.app.Validate()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := utils.GetUserClaimsFromCtx(ctx)
	body := &createMovieRequest{validator: v}

	if err := render.Bind(r, body); err != nil {
		jsonSender.BadRequest(err)
		return
	}

	movie := models.Movie{
		Name:      body.Name,
		CreatorId: userClaims.Sub,
		ListId:    userClaims.Scope,
		Status:    "pending",
	}

	_, createErr := db.NewInsert().
		Model(&movie).
		Exec(ctx)

	if createErr != nil {
		jsonSender.InternalServerError(createErr)
		return
	}

	jsonSender.Created(createMovieResponse{
		Id:        movie.Id,
		ListId:    movie.ListId,
		Name:      movie.Name,
		Status:    movie.Status,
		CreatedAt: utils.JsonEpochTime(movie.CreatedAt),
	})
}
