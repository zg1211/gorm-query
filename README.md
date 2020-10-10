# gorm-query
gorm query plugin(pagination, preload, order by)

```go
package main

import (
	"log"

	"github.com/jinzhu/gorm"
	query "github.com/zg1211/gorm-query"
)

type Car struct {
	ID     uint64
	Wheels []*Wheel
}

type Wheel struct {
	ID    uint64
	CarID uint64
}

func main() {
	db := &gorm.DB{}
	db = db.Table("cars")

	q := query.Query{
		Paginator: query.NewPaginator(10, 1),
		OrderBy:   query.NewOrderBy("id DESC"),
		Preloads: query.NewPreloads(
			map[string]func(*gorm.DB) *gorm.DB{
				"Wheels": query.NoPreloadConditions(),
			},
		),
	}

	cars := make([]*Car, 0)
	q.Find(db, &cars)

	log.Printf("%+v", q.Pagination)
}
```
