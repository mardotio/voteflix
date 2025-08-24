package models

import "github.com/uptrace/bun"

type ListUser struct {
	bun.BaseModel `bun:"table:list_users"`
	Timestamps
	UserId          string  `bun:"user_id,pk,type:uuid,notnull"`
	ListId          string  `bun:"list_id,pk,type:uuid,notnull"`
	DiscordNickname *string `bun:"discord_nickname,type:varchar(32)"`
	User            *User   `bun:"rel:belongs-to,join:user_id=id"`
	List            *List   `bun:"rel:belongs-to,join:list_id=id"`
}
