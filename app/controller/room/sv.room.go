package room

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"app/app/model"
// 	"app/app/request"
// )

// func (s *Service) Create(ctx context.Context, req request.ProductCeate) (*model.Product, bool, error) {
// 	m := model.Product{
// 		Name:        req.Name,
// 		Price:       req.Price,
// 		Description: req.Description,
// 	}
// 	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "duplicate key value") {
// 			return nil, true, errors.New("product already exists")
// 		}
// 	}
// 	return &m, false, err
// }

// func (s *Service) Update(ctx context.Context, id int64, req request.ProductUpdate) (*model.Product, bool, error) {
// 	ex, err := s.db.NewSelect().Model((*model.Product)(nil)).Where("id = ?", id).Exists(ctx)
// 	if err != nil {
// 		return nil, false, err
// 	}

// 	if !ex {
// 		return nil, true, errors.New("product not found")
// 	}

// 	m := model.Product{
// 		ID:          id,
// 		Name:        req.Name,
// 		Price:       req.Price,
// 		Description: req.Description,
// 	}

// 	m.SetUpdateNow()

// 	_, err = s.db.NewUpdate().Model(&m).
// 		Set("name = ?name").
// 		Set("price = ?price").
// 		Set("description = ?description").
// 		Set("updated_at = ?updated_at").
// 		WherePK().
// 		OmitZero().
// 		Returning("*").
// 		Exec(ctx)

// 	return &m, false, err
// }

// func (s *Service) Delete(ctx context.Context, id int64) (*model.Product, bool, error) {
// 	ex, err := s.db.NewSelect().Model((*model.Product)(nil)).Where("id = ?", id).Exists(ctx)
// 	if err != nil {
// 		return nil, false, err
// 	}

// 	if !ex {
// 		return nil, true, errors.New("product not found")
// 	}

// 	_, err = s.db.NewDelete().Model((*model.Product)(nil)).Where("id = ?", id).Exec(ctx)

// 	return nil, false, err
// }

// func (s *Service) Get(ctx context.Context, id int64) (*model.Product, error) {
// 	m := model.Product{}

// 	err := s.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
// 	return &m, err
// }

// func (s *Service) List(ctx context.Context, req request.ProductListReuest) ([]model.Product, int, error) {
// 	m := []model.Product{}

// 	var (
// 		offset = (req.Page - 1) * req.Size
// 		limit  = req.Size
// 	)

// 	query := s.db.NewSelect().Model(&m)

// 	if req.Search != "" {
// 		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
// 		query.Where("LOWER(name) Like ?", search)
// 	}

// 	count, err := query.Count(ctx)
// 	if count == 0 {
// 		return m, 0, err
// 	}

// 	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
// 	err = query.Offset(offset).Limit(limit).Order(order).Scan(ctx, &m)

// 	return m, count, err
// }

func (s *Service) Create(ctx context.Context, req request.CreateRoom) (*model.Room, bool, error) {
	m := model.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Image_url:   req.Image_url,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}
	return &m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateRoom, id request.GetByIdRoom) (*model.Room, bool, error) {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	m := &model.Room{
		ID:          id.ID,
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Image_url:   req.Image_url,
	}

	m.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		Set("description = ?description").
		Set("capacity = ?capacity").
		Set("image_url = ?image_url").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}

	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.RooomResponse, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.RooomResponse{}

	query := s.db.NewSelect().
		TableExpr("rooms as r").
		Column("r.id", "r.name", "r.description", "r.capacity", "r.image_url", "r.created_at", "r.updated_at").Where("deleted_at IS NULL").OrderExpr("r.name ASC")

		if req.Search != "" {
			searchBy := strings.ToLower(req.SearchBy)
			search := req.Search
		
			switch searchBy {
			case "name", "description":
				query = query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", searchBy), "%"+strings.ToLower(search)+"%")
			case "capacity":
				// แปลง string -> int
				if capValue, err := strconv.Atoi(search); err == nil {
					query = query.Where("r.capacity = ?", capValue)
				}
			default:
				// fallback ถ้าไม่กำหนดหรือไม่รองรับ
				query = query.Where("LOWER(r.name) LIKE ?", "%"+strings.ToLower(search)+"%")
			}
		}
		

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Order handling
	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIdRoom) (*response.RooomResponse, error) {
	m := response.RooomResponse{}

	err := s.db.NewSelect().
		TableExpr("rooms as r").
		Column("r.id", "r.name", "r.description", "r.capacity", "r.updated_at").Where("deleted_at IS NULL").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdRoom) error {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("room not found")
	}

	_, err = s.db.NewDelete().Model((*model.Room)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
