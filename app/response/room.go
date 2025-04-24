package response

type RooomResponse struct {
	ID          string `bun:"id" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	Capacity    int64  `bun:"capacity" json:"capacity"`
	Image_url   string `bun:"image_url" json:"image_url"`
	CreatedAt   string `bun:"created_at" json:"created_at"`
	UpdatedAt   string `bun:"updated_at" json:"updated_at"`
}
