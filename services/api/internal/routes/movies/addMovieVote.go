package movies

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"net/http"
	"time"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

const StatusThreshold = 2

type addMovieVoteRequest struct {
	Approve *bool `json:"approve" validate:"boolean,required"`

	validator *validator.Validate
}

func (body addMovieVoteRequest) Bind(r *http.Request) error {
	return body.validator.Struct(body)
}

type addMovieVoteResponse struct {
	MovieId   string               `json:"movieId"`
	Approved  bool                 `json:"approved"`
	CreatedAt utils.JsonEpochTime  `json:"createdAt"`
	UpdatedAt *utils.JsonEpochTime `json:"updatedAt"`
}

func (body addMovieVoteResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func createOrUpdateVote(ctx context.Context, tx bun.Tx, body addMovieVoteRequest, movie models.Movie, userClaims utils.UserJwtClaims) (vote models.Vote, created bool, updated bool, err error) {
	now := time.Now()
	vote = models.Vote{
		MovieId:    movie.Id,
		UserId:     userClaims.Sub,
		IsApproval: *body.Approve,
		Timestamps: models.Timestamps{
			CreatedAt: now,
		},
	}

	err = tx.NewInsert().
		Model(&vote).
		On("conflict (user_id, movie_id) do update").
		Set("is_approval = ?", vote.IsApproval).
		Set("updated_at = ?", now).
		Where("vote.is_approval = ?", !*body.Approve).
		Returning("*").
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return vote, false, false, nil
	}

	if err != nil {
		return vote, false, false, err
	}

	created = vote.UpdatedAt == nil
	updated = !created

	return vote, created, updated, nil
}

func updateMovieCounts(ctx context.Context, tx bun.Tx, movie models.Movie, vote models.Vote) error {
	rejectCount := movie.RejectCount
	approveCount := movie.ApproveCount

	if vote.UpdatedAt == nil {
		if vote.IsApproval {
			approveCount += 1
		} else {
			rejectCount += 1
		}
	} else {
		if vote.IsApproval {
			approveCount += 1
			rejectCount -= 1
		} else {
			approveCount -= 1
			rejectCount += 1
		}
	}

	status := movie.Status

	if rejectCount >= StatusThreshold {
		status = "rejected"
	} else if approveCount >= StatusThreshold {
		status = "accepted"
	}

	_, err := tx.NewUpdate().
		Model((*models.Movie)(nil)).
		Where("id = ?", movie.Id).
		Where("status = ?", movie.Status).
		Where("approve_count = ?", movie.ApproveCount).
		Where("reject_count = ?", movie.RejectCount).
		Set("approve_count = ?", approveCount).
		Set("reject_count = ?", rejectCount).
		Set("status = ?", status).
		Exec(ctx)

	return err
}

func (h *Handler) AddMovieVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	v := h.app.Validate()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := utils.GetUserClaimsFromCtx(ctx)
	targetMovie := middleware.GetMovieFromCtx(ctx)
	body := addMovieVoteRequest{validator: v}

	if err := render.Bind(r, &body); err != nil {
		jsonSender.BadRequest(err)
		return
	}

	if targetMovie.Status != "pending" {
		jsonSender.BadRequest(fmt.Errorf("cannot vote on movie in \"%s\" status", targetMovie.Status))
		return
	}

	tx, txErr := db.BeginTx(ctx, &sql.TxOptions{})
	defer utils.TxnRollback(&tx)

	if txErr != nil {
		jsonSender.InternalServerError(txErr)
		return
	}

	vote, didCreateVote, didUpdateVote, voteErr := createOrUpdateVote(ctx, tx, body, targetMovie, userClaims)

	if voteErr != nil {
		jsonSender.InternalServerError(voteErr)
		return
	}

	if !didCreateVote && !didUpdateVote {
		jsonSender.Conflict(fmt.Errorf("vote with same approval status already exists"))
		return
	}

	if err := updateMovieCounts(ctx, tx, targetMovie, vote); err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	if err := tx.Commit(); err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	jsonSender.Ok(addMovieVoteResponse{
		MovieId:   vote.MovieId,
		Approved:  vote.IsApproval,
		CreatedAt: utils.JsonEpochTime(vote.CreatedAt),
		UpdatedAt: (*utils.JsonEpochTime)(vote.UpdatedAt),
	})
}
