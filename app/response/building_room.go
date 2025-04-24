package response

type Building_RoomResponse struct {
    BuildingRoomID string `json:"building_room_id"` // Added field for building_room_id
    RoomID         string `json:"room_id"`
    RoomName       string `json:"room_name"`
    BuildingID     string `json:"building_id"`
    BuildingName   string `json:"building_name"`
    CreatedAt      string `json:"created_at"`
    UpdatedAt      string `json:"updated_at"`
    
}