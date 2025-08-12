package models

import "github.com/uptrace/bun"

type Movie struct {
	bun.BaseModel `bun:"table:movies"`
	PrimaryUuidId
	Timestamps
	ListId    string `bun:"list_id,type:uuid,notnull"`
	Name      string `bun:"name,type:varchar(255),notnull"`
	Status    string `bun:"status,type:varchar(15),notnull"`
	Seed      int64  `bun:"seed,type:integer,notnull"`
	CreatorId string `bun:"creator_id,type:uuid,notnull"`
	Creator   *User  `bun:"rel:belongs_to,join:creator_id=id"`
}
