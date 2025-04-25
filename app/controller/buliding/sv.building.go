package building

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateBuilding) (*model.Building, bool, error) {
	m := model.Building{
		Name: req.Name,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateBuilding, id request.GetByIdBuilding) (*model.Building, bool, error) {
	ex, err := s.db.NewSelect().Table("buildings").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Building{
		ID:          id.ID,
		Name:        req.Name,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("building already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListBuilding) ([]response.BuildingResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.BuildingResponse{}

	query := s.db.NewSelect().
		TableExpr("buildings as b").
		Column("b.id", "b.name", "b.created_at", "b.updated_at").
		Where("deleted_at IS NULL").
		OrderExpr("b.name ASC")

	// Filtering
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			searchBy := strings.ToLower(req.SearchBy)
			query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), search)
		} else {
			query = query.Where("LOWER(b.name) LIKE ?", search)
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

func (s *Service) Get(ctx context.Context, id request.GetByIdBuilding) (*response.BuildingResponse, error) {
	m := response.BuildingResponse{}

	err := s.db.NewSelect().
		TableExpr("buildings as b").
		Column("b.id", "b.name", "b.updated_at").Where("deleted_at IS NULL").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdBuilding) error {
	ex, err := s.db.NewSelect().Table("buildings").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return  err
	}
	if !ex {
		return errors.New("building not found")
	}

	_, err = s.db.NewDelete().Model((*model.Building)(nil)).Where("id = ?",id.ID).Exec(ctx)
	return err
}
