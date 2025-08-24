package users

import (
	"net/http"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type whoAmIResponse struct {
	Id          string      `json:"id"`
	DisplayName string      `json:"displayName"`
	AvatarUrl   *string     `json:"avatarUrl"`
	List        listDetails `json:"list"`
}

type listDetails struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ServerId string `json:"serverId"`
}

func (rd whoAmIResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func (h *Handler) WhoAmI(w http.ResponseWriter, r *http.Request) {
	db := h.app.Db()
	ctx := r.Context()
	jsonSender := utils.NewJsonSender(w, r)

	userClaims := utils.GetUserClaimsFromCtx(ctx)

	user := models.ListUser{}

	userQuery := db.NewSelect().
		Model(&user).
		Where("list_user.id = ?", userClaims.Sub).
		Relation("User").
		Relation("List")

	if err := userQuery.Scan(ctx); err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	response := whoAmIResponse{
		Id: user.Id,
		List: listDetails{
			Id:       user.List.Id,
			Name:     user.List.Name,
			ServerId: user.List.DiscordServerId,
		},
		AvatarUrl:   utils.GetAvatarUrl(*user.User),
		DisplayName: user.User.DiscordUsername,
	}

	if user.DiscordNickname != nil {
		response.DisplayName = *user.DiscordNickname
	}

	jsonSender.Ok(response)
}
