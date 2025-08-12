package models

import "time"

type PrimaryUuidId struct {
	Id string `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
}

type Timestamps struct {
	CreatedAt time.Time  `bun:"created_at,type:timestamp,default:current_timestamp"`
	UpdatedAt *time.Time `bun:"updated_at,type:timestamp"`
}
