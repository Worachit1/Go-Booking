package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Building_Room(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	building_room := router.Group("")
	{
		building_room.POST("/create", ctl.Building_RoomCtl.Create)
		building_room.PATCH("/:id", ctl.Building_RoomCtl.Update)
		building_room.GET("/list", ctl.Building_RoomCtl.List)
		building_room.GET("/:id", ctl.Building_RoomCtl.Get)
		building_room.DELETE("/:id", ctl.Building_RoomCtl.Delete)
	}
}