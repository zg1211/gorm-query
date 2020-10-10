package query

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Query struct {
	Preloads
	Paginator
	OrderBy

	Pagination *Pagination
}

func (q *Query) First(db *gorm.DB, data interface{}) error {
	query := db

	if q.OrderBy != nil {
		query = q.Order(query)
	}

	if q.Preloads != nil {
		query = q.Load(query)
	}

	return query.First(data).Error
}

func (q *Query) Find(db *gorm.DB, data interface{}) error {
	query := db

	if q.OrderBy != nil {
		query = q.Order(query)
	}

	if q.Preloads != nil {
		query = q.Load(query)
	}

	if q.Paginator != nil {
		pagination, err := q.Page(query, data)
		if err != nil {
			return err
		}

		q.Pagination = pagination

		return nil
	} else {
		if err := query.Find(data).Error; err != nil {
			return err
		}
	}

	return nil
}

type OrderBy interface {
	Order(*gorm.DB) *gorm.DB
}

type orderBy struct {
	Expr interface{}
}

func NewOrderBy(value string, args ...interface{}) *orderBy {
	return &orderBy{
		Expr: gorm.Expr(value, args...),
	}
}

func (o *orderBy) Order(db *gorm.DB) *gorm.DB {
	return db.Order(o.Expr)
}

type Preloads interface {
	Load(db *gorm.DB) *gorm.DB
}

type preloads map[string]func(db *gorm.DB) *gorm.DB

func NewPreloads(p map[string]func(*gorm.DB) *gorm.DB) preloads {
	return preloads(p)
}

func NoPreloadConditions() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func (p preloads) Load(db *gorm.DB) *gorm.DB {
	for k, v := range p {
		db = db.Preload(k, v)
	}

	return db
}

type Paginator interface {
	Page(db *gorm.DB, data interface{}) (*Pagination, error)
}

type paginator struct {
	per, page uint32
}

type Pagination struct {
	Per        uint32
	Page       uint32
	TotalPage  uint32
	TotalCount uint32
	HasMore    bool
}

func NewPaginator(per, page uint32) *paginator {
	return &paginator{
		per:  per,
		page: page,
	}
}

func (p *paginator) Page(db *gorm.DB, data interface{}) (*Pagination, error) {
	if p.per == 0 || p.page == 0 {
		return nil, fmt.Errorf("Invalid Per: %d, Page: %d", p.per, p.page)
	}

	var totalCount, totalPage uint32
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := db.Offset(p.per * (p.page - 1)).Limit(p.per).Find(data).Error; err != nil {
		return nil, err
	}

	if totalCount == 0 {
		totalPage = 0
	} else {
		totalPage = (totalCount-1)/p.per + 1
	}

	return &Pagination{
		Per:        p.per,
		Page:       p.page,
		TotalCount: totalCount,
		TotalPage:  totalPage,
		HasMore:    totalPage > p.page,
	}, nil
}
