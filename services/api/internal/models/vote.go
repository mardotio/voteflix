package models

import "github.com/uptrace/bun"

type Vote struct {
	bun.BaseModel `bun:"table:votes"`
	Timestamps
	MovieId    string `bun:"movie_id,pk,type:uuid,notnull"`
	UserId     string `bun:"user_id,pk,type:uuid,notnull"`
	IsApproval bool   `bun:"is_approval,type:boolean,notnull"`
	Movie      *Movie `bun:"rel:belongs-to,join:movie_id=id"`
	User       *User  `bun:"rel:belongs-to,join:user_id=id"`
}
