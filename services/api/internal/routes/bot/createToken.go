package bot

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
	"voteflix/api/internal/utils"
)

type createTokenRequest struct {
	Sub       string  `json:"sub" validate:"required,max=100"`
	Server    string  `json:"server" validate:"required,max=100"`
	Username  string  `json:"username" validate:"required,max=100"`
	Avatar    *string `json:"avatar" validate:"omitempty,max=100"`
	Nickname  *string `json:"nickname" validate:"omitempty,max=100"`
	validator *validator.Validate
}

type createTokenResponse struct {
	Token     string              `json:"token"`
	ExpiresAt utils.JsonEpochTime `json:"expiresAt"`
}

func (body createTokenRequest) Bind(r *http.Request) error {
	return body.validator.Struct(body)
}

func (rd createTokenResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	v := h.app.Validate()
	body := &createTokenRequest{validator: v}
	jsonSender := utils.NewJsonSender(w, r)

	if err := render.Bind(r, body); err != nil {
		jsonSender.BadRequest(err)
		return
	}

	claims := utils.BotJwtClaims{
		Sub:      body.Sub,
		Server:   body.Server,
		Username: body.Username,
		Avatar:   body.Avatar,
		Nickname: body.Nickname,
	}
	token, tokenString, _ := h.app.BotJwtAuth().Encode(claims.ToClaimsMap(time.Minute * 2))
	jsonSender.Created(createTokenResponse{Token: tokenString, ExpiresAt: utils.JsonEpochTime(token.Expiration())})
}
