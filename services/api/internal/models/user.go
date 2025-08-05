package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	Id            string    `bun:"type:uuid,pk,default:gen_random_uuid()"`
	CreatedAt     time.Time `bun:"default:current_timestamp"`
}
