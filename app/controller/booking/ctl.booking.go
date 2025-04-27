package booking

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)


func (ctl *Controller) Create(ctx *gin.Context) {
	req := request.CreateBooking{}
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Create(ctx, req)
	if err != nil {
		ms := "Internal Server Error"
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIdBooking{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}
	body := request.UpdateBooking{}
	if err := ctx.Bind(&body); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := "Internal Server Error"
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListBooking{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0{
		req.Page = 1
	}
	if req.Page == 0 {
		req.Page = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "asc"
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil{
		logger.Errf(err.Error())
		response.InternalError(ctx,err.Error())
		return
	}
	response.SuccessWithPaginate(ctx,data,req.Size,req.Page,total)

}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIdBooking{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx,ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx,err.Error())
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) GetByRoomId(ctx *gin.Context) {
    // สมมติว่า roomID เป็น string ที่ได้จาก URL หรือการส่งข้อมูล
    roomID := ctx.Param("id") // หรือวิธีที่คุณได้ค่า roomID

    // สร้างอ็อบเจ็กต์ request.GetByRoomIdBooking และตั้งค่า RoomID
    req := request.GetByRoomIdBooking{
        RoomID: roomID, // ตั้งค่า RoomID ให้กับตัวแปรที่รับมา
        // อาจจะมีค่าต่างๆ อื่นๆ เช่น Page, Size ที่จะตั้งค่าเช่นกัน
        Page: 1,  // ตัวอย่างค่า Page
        Size: 10, // ตัวอย่างค่า Size
    }

    // ส่งอ็อบเจ็กต์ req ไปยังฟังก์ชัน Service.GetByRoomId
    data, total, err := ctl.Service.GetByRoomId(ctx, req)

    if err != nil {
        logger.Errf(err.Error())
        response.InternalError(ctx, err.Error())
        return
    }
    response.SuccessWithPaginate(ctx, data, 0, 0, total)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIdBooking{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	err := ctl.Service.Delete(ctx,ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx,err.Error())
		return
	}
	response.Success(ctx, nil)
}