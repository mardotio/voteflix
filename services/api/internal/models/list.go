package models

import "github.com/uptrace/bun"

type List struct {
	bun.BaseModel `bun:"table:lists"`
	PrimaryUuidId
	Timestamps
	Name            string  `bun:"name,type:varchar(100),notnull"`
	DiscordServerId string  `bun:"discord_server_id,type:varchar(100),notnull,unique"`
	DiscordAvatarId *string `bun:"discord_avatar_id,type:varchar(100),notnull"`
	CreatorId       string  `bun:"creator_id,type:uuid"`
	Creator         *User   `bun:"rel:belongs-to,join:creator_id=id"`
}
