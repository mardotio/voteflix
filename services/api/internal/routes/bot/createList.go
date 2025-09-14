package bot

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"net/http"
	"time"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type createListRequest struct {
	DiscordUserId         string  `json:"discordUserId" validate:"required,max=100"`
	DiscordUsername       string  `json:"discordUsername" validate:"required,max=100"`
	DiscordServerId       string  `json:"discordServerId" validate:"required,max=100"`
	DiscordServerName     string  `json:"discordServerName" validate:"required,max=100"`
	DiscordServerAvatarId *string `json:"discordServerAvatarId" validate:"omitempty,max=100"`
	DiscordNickname       *string `json:"discordNickname" validate:"omitempty,max=100"`
	DiscordAvatarId       *string `json:"discordAvatarId" validate:"omitempty,max=100"`
	validator             *validator.Validate
}

type createListResponse struct {
	Id string `json:"id"`
}

func (body createListRequest) Bind(r *http.Request) error {
	return body.validator.Struct(body)
}

func (rd createListResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func getUserId(ctx context.Context, tx bun.Tx, body *createListRequest) (string, error) {
	now := time.Now()

	newUser := models.User{
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: &now,
		},
		DiscordId:       body.DiscordUserId,
		DiscordUsername: body.DiscordUsername,
		DiscordAvatarId: body.DiscordAvatarId,
	}

	_, newUserError := tx.NewInsert().
		Model(&newUser).
		On("conflict (discord_id) do update").
		Set("updated_at = EXCLUDED.updated_at").
		Set("discord_username = EXCLUDED.discord_username").
		Set("discord_avatar_id = EXCLUDED.discord_avatar_id").
		Exec(ctx)

	if newUserError != nil {
		return "", newUserError
	}

	return newUser.Id, nil
}

func createList(ctx context.Context, tx bun.Tx, body *createListRequest, userId string) (*models.List, error) {
	list := models.List{
		CreatorId:       userId,
		Name:            body.DiscordServerName,
		DiscordServerId: body.DiscordServerId,
		DiscordAvatarId: body.DiscordServerAvatarId,
	}

	_, err := tx.NewInsert().Model(&list).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &list, nil
}

func createListUser(ctx context.Context, tx bun.Tx, body *createListRequest, userId string, listId string) (*models.ListUser, error) {
	listUser := models.ListUser{
		UserId:          userId,
		ListId:          listId,
		DiscordNickname: body.DiscordNickname,
	}

	_, err := tx.NewInsert().Model(&listUser).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &listUser, nil
}

func (h *Handler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	v := h.app.Validate()
	body := &createListRequest{validator: v}
	jsonSender := utils.NewJsonSender(w, r)

	if err := render.Bind(r, body); err != nil {
		jsonSender.BadRequest(err)
		return
	}

	existingListQuery := db.NewSelect().
		Model((*models.List)(nil)).
		Where("discord_server_id = ?", body.DiscordServerId)

	if exists, err := existingListQuery.Exists(ctx); err != nil {
		jsonSender.InternalServerError(err)
		return
	} else if exists {
		jsonSender.Conflict(fmt.Errorf("list already exists"))
		return
	}

	tx, txErr := db.BeginTx(ctx, &sql.TxOptions{})
	defer utils.TxnRollback(&tx)

	if txErr != nil {
		jsonSender.InternalServerError(txErr)
		return
	}

	userId, userIdErr := getUserId(ctx, tx, body)

	if userIdErr != nil {
		jsonSender.InternalServerError(userIdErr)
		return
	}

	list, listErr := createList(ctx, tx, body, userId)

	if listErr != nil {
		jsonSender.InternalServerError(listErr)
		return
	}

	_, listUserErr := createListUser(ctx, tx, body, userId, list.Id)

	if listUserErr != nil {
		jsonSender.InternalServerError(listUserErr)
		return
	}

	if err := tx.Commit(); err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	jsonSender.Created(createListResponse{Id: list.Id})
}
