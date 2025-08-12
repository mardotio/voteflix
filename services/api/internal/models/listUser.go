package models

import "github.com/uptrace/bun"

type ListUser struct {
	bun.BaseModel `bun:"table:list_users"`
	PrimaryUuidId
	Timestamps
	UserId          string  `bun:"user_id,type:uuid,notnull"`
	ListId          string  `bun:"list_id,type:uuid,notnull"`
	DiscordNickname *string `bun:"discord_nickname,type:varchar(255)"`
	User            *User   `bun:"rel:belongs_to,join:user_id=id"`
	List            *List   `bun:"rel:belongs_to,join:list_id=id"`
}
