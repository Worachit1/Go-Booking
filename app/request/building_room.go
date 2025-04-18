package request

type CreateBuilding_Room struct {
	RoomID     string `form:"room_id" binding:"required"`
	BuildingID string `form:"building_id" binding:"required"`
}

type UpdateBuilding_Room struct {
	CreateBuilding_Room
}

type ListBuilding_Room struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type GetByIdBuilding_Room struct {
	ID string `uri:"id" binding:"required"`
}
