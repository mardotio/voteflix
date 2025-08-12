package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	PrimaryUuidId
	Timestamps
	DiscordId       string  `bun:"discord_id,type:varchar(100),notnull,unique"`
	DiscordUsername string  `bun:"discord_username,type:varchar(100),notnull"`
	DiscordAvatarId *string `bun:"discord_avatar_id,type:varchar(100)"`
}
