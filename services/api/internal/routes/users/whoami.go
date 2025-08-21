package users

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
	"voteflix/api/internal/utils"
)

type whoAmIResponse struct {
	UserId string `json:"userId"`
}

func (rd whoAmIResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) WhoAmI(w http.ResponseWriter, r *http.Request) {
	userClaims := utils.GetUserClaimsFromCtx(r.Context())
	render.Status(r, http.StatusOK)
	err := render.Render(w, r, whoAmIResponse{UserId: userClaims.Sub})

	if err != nil {
		log.Println(err)
	}
}
