package models

import "github.com/uptrace/bun"

type Rating struct {
	bun.BaseModel `bun:"table:ratings"`
	Timestamps
	MovieId    string `bun:"movie_id,pk,type:uuid,notnull"`
	ListUserId string `bun:"list_user_id,pk,type:uuid,notnull"`
	Rating     int64  `bun:"rating,type:smallint,notnull"`
	Movie      *Movie `bun:"rel:belongs_to,join:movie_id=id"`
	List       *List  `bun:"rel:belongs_to,join:list_id=id"`
}
