package model

import (
	"app/app/enum"
	"time"

	"github.com/uptrace/bun"
)

type Booking struct {
	bun.BaseModel `bun:"table:bookings"`

	ID          string             `bun:",pk,type:uuid,default:gen_random_uuid()"`
	User        *User              `bun:"rel:belongs-to,join:user_id=id"`
	UserID      string             `bun:"user_id,notnull"`
	Room        *Room              `bun:"rel:belongs-to,join:room_id=id"`
	RoomID      string             `bun:"room_id,notnull"`
	Title       string             `bun:"title,notnull"`
	Description string             `bun:"description,notnull"`
	StartTime   time.Time          `bun:"start_time,notnull"`
	EndTime     time.Time          `bun:"end_time,notnull"`
	Status      enum.BookingStatus `bun:"status,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
