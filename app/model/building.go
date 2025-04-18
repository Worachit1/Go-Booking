package model

import (
	"github.com/uptrace/bun"
)

type Building struct {
	bun.BaseModel `bun:"table:buildings"`

	ID   string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name string `bun:"name,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
