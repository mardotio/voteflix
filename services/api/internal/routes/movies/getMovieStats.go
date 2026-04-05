package movies

import (
	"context"
	"net/http"
	"voteflix/api/internal/middleware"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"

	"github.com/uptrace/bun"
)

type getMovieStatsResponse struct {
	Pending  int64 `json:"pending"`
	Approved int64 `json:"approved"`
	Watched  int64 `json:"watched"`
	Rejected int64 `json:"rejected"`
	Total    int64 `json:"total"`
}

func (body getMovieStatsResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

type statusCounts struct {
	Status string
	Count  int64
}

func getStats(db *bun.DB, ctx context.Context, list string) (*getMovieStatsResponse, error) {
	var stats []statusCounts

	err := db.NewSelect().
		Model((*models.Movie)(nil)).
		Column("status").
		ColumnExpr("count(*) as count").
		Where("list_id = ?", list).
		Group("status").
		Scan(ctx, &stats)

	if err != nil {
		return nil, err
	}

	stdStats := getMovieStatsResponse{
		Pending:  0,
		Approved: 0,
		Watched:  0,
		Rejected: 0,
	}

	if len(stats) == 0 {
		return &stdStats, nil
	}

	for _, stat := range stats {
		if stat.Status == "pending" {
			stdStats.Pending = stat.Count
		}
		if stat.Status == "approved" {
			stdStats.Approved = stat.Count
		}
		if stat.Status == "watched" {
			stdStats.Watched = stat.Count
		}
		if stat.Status == "rejected" {
			stdStats.Rejected = stat.Count
		}

		stdStats.Total += stat.Count
	}

	return &stdStats, nil
}

func (h *Handler) GetMovieStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)
	userClaims := middleware.GetUserClaimsFromCtx(ctx)

	stats, err := getStats(db, ctx, userClaims.Scope)

	if err != nil {
		jsonSender.InternalServerError(err)
	}

	jsonSender.Ok(stats)
}
