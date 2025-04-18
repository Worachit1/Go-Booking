package booking

import "github.com/uptrace/bun"

type Controller struct {
	Name    string
	Service *Service
}

func NewController(db *bun.DB) *Controller {
	return &Controller{
		Name:    `booking-ctl`,
		Service: NewService(db),
	}
}

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{
		db: db,
	}
}
