package booking

import (
	"app/app/enum"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func parseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic("invalid time format, expected RFC3339")
	}
	return parsedTime
}

func (s *Service) Create(ctx context.Context, req request.CreateBooking) (*model.Booking, bool, error) {
	m := model.Booking{
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		Status:      enum.BookingStatus("Pending"),
	}

	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("booking already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateBooking, id request.GetByIdBooking) (*model.Booking, bool, error) {
	ex, err := s.db.NewSelect().Table("bookings").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Booking{
		ID:          id.ID,
		UserID:      req.UserID,
		RoomID:      req.RoomID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		Status:      enum.BookingStatus(req.Status),
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(m).
		Set("room_id = ?room_id").
		Set("description = ?description").
		Set("start_time = ?start_time").
		Set("end_time = ?end_time").
		Set("status = ?status").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("booking already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListBooking) ([]response.BookingResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.BookingResponse{}

	// สร้าง query base ที่ใช้ร่วมกัน
	baseQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.created_at as created_at").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Where("b.deleted_at IS NULL").
		OrderExpr("b.start_time DESC")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			baseQuery = baseQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			baseQuery = baseQuery.Where("LOWER(b.title) LIKE ?", search) // เปลี่ยนจาก name เป็น title
		}
	}

	// ใช้ baseQuery สำหรับการนับจำนวนข้อมูล
	countQuery := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("COUNT(*)").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Where("b.deleted_at IS NULL")

	// ใช้การ filter เดียวกับการดึงข้อมูลหลัก
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			countQuery = countQuery.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", searchBy), search)
		} else {
			countQuery = countQuery.Where("LOWER(b.title) LIKE ?", search)
		}
	}

	// คำนวณจำนวนทั้งหมด
	var count int
	err := countQuery.Scan(ctx, &count)
	if err != nil {
		return nil, 0, err
	}

	// Final query สำหรับดึงข้อมูล
	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)
	err = baseQuery.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdBooking) (*response.BookingResponse, error) {
	m := response.BookingResponse{}

	err := s.db.NewSelect().
		TableExpr("bookings as b").
		ColumnExpr("b.id as id").
		ColumnExpr("u.id as user_id").
		ColumnExpr("u.first_name as user_name").
		ColumnExpr("u.last_name as user_lastname").
		ColumnExpr("r.id as room_id").
		ColumnExpr("r.name as room_name").
		ColumnExpr("b.title as title").
		ColumnExpr("b.description as description").
		ColumnExpr("b.start_time as start_time").
		ColumnExpr("b.end_time as end_time").
		ColumnExpr("b.status as status").
		ColumnExpr("b.updated_at as updated_at").
		Join("JOIN users as u ON b.user_id::uuid = u.id").
		Join("JOIN rooms as r ON b.room_id::uuid = r.id").
		Where("b.deleted_at IS NULL").
		OrderExpr("b.start_time DESC").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) GetByRoomId(ctx context.Context, req request.GetByRoomIdBooking) ([]response.BookingResponse, int, error) {
    // เพิ่มการรับค่า Page และ Size จาก req
    offset := (req.Page - 1) * req.Size
    m := []response.BookingResponse{}

    query := s.db.NewSelect().
        TableExpr("bookings as b").
        ColumnExpr("b.id as id").
        ColumnExpr("u.id as user_id").
        ColumnExpr("u.first_name as user_name").
        ColumnExpr("u.last_name as user_lastname").
        ColumnExpr("r.id as room_id").
        ColumnExpr("r.name as room_name").
        ColumnExpr("b.title as title").
        ColumnExpr("b.description as description").
        ColumnExpr("b.start_time as start_time").
        ColumnExpr("b.end_time as end_time").
        ColumnExpr("b.status as status").
        ColumnExpr("b.created_at as created_at").
        ColumnExpr("b.updated_at as updated_at").
        Join("JOIN users as u ON b.user_id::uuid = u.id").
        Join("JOIN rooms as r ON b.room_id::uuid = r.id").
        Where("b.deleted_at IS NULL").
        Where("r.id = ?", req.RoomID). // ตรวจสอบว่า RoomID ถูกต้องและเป็น UUID
        OrderExpr("start_time ASC")

    count, err := query.Count(ctx)
    if err != nil {
        return nil, 0, err
    }

    err = query.OrderExpr("start_time ASC").Limit(req.Size).Offset(offset).Scan(ctx, &m)
    if err != nil {
        return nil, 0, err
    }

    return m, count, nil
}


func (s *Service) Delete(ctx context.Context, id request.GetByIdBooking) error {
	ex, err := s.db.NewSelect().Table("bookings").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("booking not found")
	}

	_, err = s.db.NewDelete().Model((*model.Booking)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
