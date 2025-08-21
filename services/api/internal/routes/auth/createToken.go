package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"net/http"
	"time"
	"voteflix/api/internal/models"
	"voteflix/api/internal/utils"
)

type createTokenResponse struct {
	Token     string              `json:"token"`
	ExpiresAt utils.JsonEpochTime `json:"expiresAt"`
}

func (rd createTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getList(ctx context.Context, db *bun.DB, serverId string) (string, error) {
	list := models.List{}

	query := db.NewSelect().Model(&list).Where("discord_server_id = ?", serverId)

	if err := query.Scan(ctx); err != nil {
		return "", err
	}

	return list.Id, nil
}

func getUserId(ctx context.Context, tx bun.Tx, claims utils.BotJwtClaims) (string, error) {
	now := time.Now()

	existingUser := models.User{}

	existingUserErr := tx.NewSelect().
		Model(&existingUser).
		Where("discord_id = ?", claims.Sub).
		Scan(ctx)

	if existingUserErr != nil && !errors.Is(existingUserErr, sql.ErrNoRows) {
		return "", existingUserErr
	}

	if existingUserErr == nil {
		shouldUpdate := false

		if existingUser.DiscordUsername != claims.Username {
			shouldUpdate = true
		} else if !utils.IsEqual(existingUser.DiscordAvatarId, claims.Avatar) {
			shouldUpdate = true
		}

		if !shouldUpdate {
			return existingUser.Id, nil
		}

		_, updateErr := tx.NewUpdate().
			Model(&existingUser).
			Where("id = ?", existingUser.Id).
			Set("discord_username = ?", claims.Username).
			Set("discord_avatar_id = ?", claims.Avatar).
			Set("updated_at = ?", now).
			Exec(ctx)

		if updateErr != nil {
			return "", updateErr
		}

		return existingUser.Id, nil
	}

	newUser := models.User{
		Timestamps: models.Timestamps{
			CreatedAt: now,
		},
		DiscordId:       claims.Sub,
		DiscordUsername: claims.Username,
		DiscordAvatarId: claims.Avatar,
	}

	_, newUserError := tx.NewInsert().
		Model(&newUser).
		Exec(ctx)

	if newUserError != nil {
		return "", newUserError
	}

	return newUser.Id, nil
}

func getListUser(ctx context.Context, tx bun.Tx, claims utils.BotJwtClaims, userId string, listId string) (models.ListUser, error) {
	now := time.Now()

	existingUser := models.ListUser{}

	existingUserErr := tx.NewSelect().
		Model(&existingUser).
		Where("user_id = ?", userId).
		Where("list_id = ?", listId).
		Scan(ctx)

	if existingUserErr != nil && !errors.Is(existingUserErr, sql.ErrNoRows) {
		return existingUser, existingUserErr
	}

	if existingUserErr == nil {
		shouldUpdate := !utils.IsEqual(existingUser.DiscordNickname, claims.Nickname)

		if !shouldUpdate {
			return existingUser, nil
		}

		_, updateErr := tx.NewUpdate().
			Model(&existingUser).
			Where("id = ?", existingUser.Id).
			Set("discord_nickname = ?", claims.Nickname).
			Set("updated_at = ?", now).
			Exec(ctx)

		if updateErr != nil {
			return existingUser, updateErr
		}

		return existingUser, nil
	}

	newUser := models.ListUser{
		Timestamps: models.Timestamps{
			CreatedAt: now,
		},
		DiscordNickname: claims.Nickname,
		UserId:          userId,
		ListId:          listId,
	}

	_, newUserError := tx.NewInsert().
		Model(&newUser).
		Exec(ctx)

	if newUserError != nil {
		return newUser, newUserError
	}

	return newUser, nil
}

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := h.app.Db()
	jsonSender := utils.NewJsonSender(w, r)
	botClaims := utils.GetBotClaimsFromCtx(ctx)

	serverListId, serverListErr := getList(ctx, db, botClaims.Server)

	if serverListErr != nil {
		jsonSender.NotFound(fmt.Errorf("could not find list assocaited with server %s", botClaims.Server))
		return
	}

	tx, txErr := db.BeginTx(ctx, &sql.TxOptions{})
	defer utils.TxnRollback(&tx)

	if txErr != nil {
		jsonSender.InternalServerError(txErr)
		return
	}

	userId, userIdErr := getUserId(ctx, tx, botClaims)

	if userIdErr != nil {
		jsonSender.InternalServerError(userIdErr)
		return
	}

	listUser, listUserErr := getListUser(ctx, tx, botClaims, userId, serverListId)

	if listUserErr != nil {
		jsonSender.InternalServerError(listUserErr)
		return
	}

	if err := tx.Commit(); err != nil {
		jsonSender.InternalServerError(err)
		return
	}

	claims := utils.UserJwtClaims{Sub: listUser.Id, List: serverListId}
	token, tokenString, _ := h.app.JwtAuth().Encode(claims.ToClaimsMap(time.Hour * 1))
	jsonSender.Created(createTokenResponse{Token: tokenString, ExpiresAt: utils.JsonEpochTime(token.Expiration())})
}
