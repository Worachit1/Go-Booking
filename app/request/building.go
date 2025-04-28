package request

type CreateBuilding struct {
	Name string  `json:"name" form:"name"` 
}

type UpdateBuilding struct{
	CreateBuilding
}

type ListBuilding struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type GetByIdBuilding struct{
	ID string `uri:"id" binding:"required"`
}