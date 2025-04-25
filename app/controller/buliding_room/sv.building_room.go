package building_room

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateBuilding_Room) (*model.Building_Room, bool, error) {
	
	exists, err := s.db.NewSelect().
	Model((*model.Building_Room)(nil)).
	Where("room_id = ?", req.RoomID).
	Exists(ctx)
	if err != nil {
		return nil, false, err 
	}
	if exists {
		return nil, true, errors.New("room_id นี้มีอยู่แล้วใน building_rooms")
	}
	
	m := &model.Building_Room{
		RoomID: req.RoomID,
		BuildingID: req.BuildingID,
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, false, err // ถ้าเกิด error ในการ insert
	}
	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateBuilding_Room, id request.GetByIdBuilding_Room) (*model.Building_Room, bool, error) {
	ex, err := s.db.NewSelect().
		Model((*model.Building_Room)(nil)).
		Where("room_id = ?", req.RoomID).
		Where("id != ?", id.ID).
		Exists(ctx)

	if err != nil {
		return nil, false, err
	}
	if ex {
		return nil, true, errors.New("room_id นี้มีอยู่แล้ว")
	}

	m := &model.Building_Room{
		ID:          id.ID,
		RoomID:      req.RoomID,
		BuildingID:  req.BuildingID,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(m).
		Set("room_id = ?room_id").
		Set("building_id = ?building_id").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building_room already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListBuilding_Room) ([]response.Building_RoomResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.Building_RoomResponse{}

	query := s.db.NewSelect().
		TableExpr("building_rooms AS br").
		ColumnExpr("br.id AS building_room_id").
		ColumnExpr("r.id AS room_id").
		ColumnExpr("r.name AS room_name").
		ColumnExpr("b.id AS building_id").
		ColumnExpr("b.name AS building_name").
		ColumnExpr("br.created_at AS created_at").
		ColumnExpr("br.updated_at AS updated_at").
		Join("JOIN rooms AS r ON br.room_id::uuid = r.id").
		Join("JOIN buildings AS b ON br.building_id::uuid = b.id").
		Where("br.deleted_at IS NULL").
		OrderExpr("r.name ASC")


	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(r.name) LIKE ? OR LOWER(bd.name) LIKE ?", search, search)
		}
	}

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order handling
	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdBuilding_Room) (*response.Building_RoomResponse, error) {
	m := response.Building_RoomResponse{}

	err := s.db.NewSelect().
	TableExpr("building_rooms AS br").
	ColumnExpr("br.id AS building_room_id").
	ColumnExpr("r.id AS room_id").
	ColumnExpr("r.name AS room_name").
	ColumnExpr("b.id AS building_id").
	ColumnExpr("b.name AS building_name").
	Join("JOIN rooms AS r ON br.room_id::uuid = r.id"). // 👈 cast ถ้า type ไม่ตรง
	Join("JOIN buildings AS b ON br.building_id::uuid = b.id").
	Where("br.deleted_at IS NULL").
	Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdBuilding_Room) error {
	ex, err := s.db.NewSelect().Table("building_rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return  err
	}
	if !ex {
		return errors.New("building_room not found")
	}

	_, err = s.db.NewDelete().Model((*model.Building_Room)(nil)).Where("id = ?",id.ID).Exec(ctx)
	return err
}
