package model

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName string `bun:"first_name,notnull"`
	LastName  string `bun:"last_name,notnull"`
	Email     string `bun:"email,unique,notnull"`
	Password  string `bun:"password,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
