package models

import "github.com/uptrace/bun"

type Vote struct {
	bun.BaseModel `bun:"table:vote"`
	Timestamps
	MovieId    string `bun:"movie_id,pk,type:uuid,notnull"`
	ListUserId string `bun:"list_user_id,pk,type:uuid,notnull"`
	IsApproval bool   `bun:"is_approval,type:boolean,notnull"`
	Movie      *Movie `bun:"rel:belongs_to,join:movie_id=id"`
	List       *List  `bun:"rel:belongs_to,join:list_id=id"`
}
