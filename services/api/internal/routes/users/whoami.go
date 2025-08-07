package users

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type whoAmIResponse struct {
	UserId string `json:"userId"`
}

func (rd whoAmIResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) WhoAmI(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	render.Status(r, http.StatusOK)
	err := render.Render(w, r, whoAmIResponse{UserId: claims["sub"].(string)})

	if err != nil {
		log.Println(err)
	}
}
