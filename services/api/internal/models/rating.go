package models

import "github.com/uptrace/bun"

type Rating struct {
	bun.BaseModel `bun:"table:ratings"`
	Timestamps
	MovieId string `bun:"movie_id,pk,type:uuid,notnull"`
	UserId  string `bun:"user_id,pk,type:uuid,notnull"`
	Rating  int64  `bun:"rating,type:smallint,notnull"`
	Movie   *Movie `bun:"rel:belongs-to,join:movie_id=id"`
	User    *User  `bun:"rel:belongs-to,join:user_id=id"`
}
