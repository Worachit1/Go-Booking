package request

type CreateBooking struct {
	UserID      string `json:"user_id"`
	RoomID      string `json:"room_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
}

type UpdateBooking struct {
	CreateBooking
}

type ListBooking struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	RoomID   string `form:"room_id"`
}

type GetByIdBooking struct {
	ID string `uri:"id" binding:"required"`
}

type GetByRoomIdBooking struct {
	ID string `uri:"id" binding:"required"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	RoomID   string `form:"room_id"`
}