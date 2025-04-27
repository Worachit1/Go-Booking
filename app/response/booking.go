package response

type BookingResponse struct {
	ID           string `json:"id" bun:"id"`
	UserID       string `json:"user_id" bun:"user_id"`
	UserName     string `json:"user_name" bun:"user_name"`
	UserLastName string `json:"user_lastname" bun:"user_lastname"`
	RoomID       string `json:"room_id" bun:"room_id"`
	RoomName     string `json:"room_name" bun:"room_name"`
	Title        string `json:"title" bun:"title"`
	Description  string `json:"description" bun:"description"`
	StartTime    string `json:"start_time" bun:"start_time"`
	EndTime      string `json:"end_time" bun:"end_time"`
	Status       string `json:"status" bun:"status"`
	CreatedAt    int64 `json:"created_at" bun:"created_at"`
	UpdatedAt    int64 `json:"updated_at" bun:"updated_at"`
}
