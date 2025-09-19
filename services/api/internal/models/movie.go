package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Movie struct {
	bun.BaseModel `bun:"table:movies"`
	PrimaryUuidId
	Timestamps
	ListId       string     `bun:"list_id,type:uuid,notnull"`
	Name         string     `bun:"name,type:varchar(255),notnull"`
	Status       string     `bun:"status,type:varchar(15),notnull"` //pending, accepted, rejected, watched
	Seed         int64      `bun:"seed,type:integer,default:gen_random_movie_seed()"`
	ApproveCount int64      `bun:"approve_count,type:integer,default:0"`
	RejectCount  int64      `bun:"reject_count,type:integer,default:0"`
	CreatorId    string     `bun:"creator_id,type:uuid,notnull"`
	WatchedAt    *time.Time `bun:"watched_at,type:timestamp"`
	Creator      *ListUser  `bun:"rel:belongs-to,join:creator_id=user_id,join:list_id=list_id"`
	User         *User      `bun:"rel:belongs-to,join:creator_id=id"`
}
