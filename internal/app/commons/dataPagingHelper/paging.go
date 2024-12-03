// Datapaging is a helper to
// create a pagination for data retrieve from repository layer
// Author : Gilang Prambudi - Jul, 18 2020
package datapaging

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Datapaging struct {
	Limit int
	Page  int

	//OrderBy define the order property of the row, use it with format [field] [asc/desc].
	//Example []string{"distance desc"}
	OrderBy []string

	//FilterColumn specify column as filter parameter
	FilterColumn string
	//FilterValue specify value of the column to be filtered
	FilterValue string

	DateLatest   *time.Time
	DateEarliest *time.Time
}

// New will return a new pagination object specified with pagination,
// limit and order by field
func New(limit, page int, orderBy []string) Datapaging {
	return Datapaging{
		Limit:   limit,
		Page:    page,
		OrderBy: orderBy,
	}
}

// NoPagination return empty pagination
func NoPagination() Datapaging {
	return Datapaging{}
}

// IsNil check if the pagination object is empty
func (pagination *Datapaging) IsNil() bool {
	if !pagination.WithLimit() && !pagination.WithPageOffset() && !pagination.WithOrderBy() {
		return true
	}
	return false
}

func (pagination *Datapaging) GetOffset() int {
	return (pagination.Page - 1) * pagination.Limit
}

func (pagination *Datapaging) WithLimit() bool {
	if pagination.Limit != 0 {
		return true
	}
	return false
}

func (pagination *Datapaging) WithPageOffset() bool {
	if pagination.Page != 0 {
		return true
	}
	return false
}

func (pagination *Datapaging) WithOrderBy() bool {
	if len(pagination.OrderBy) > 0 {
		return true
	}
	return false
}

func (pagination Datapaging) Between(earliestTime, latestTime *time.Time) Datapaging {
	pagination.DateEarliest = earliestTime
	pagination.DateLatest = latestTime
	return pagination
}

func (pagination *Datapaging) WithDateBetween() bool {
	if pagination.DateEarliest != nil && pagination.DateLatest != nil {
		return true
	}

	return false
}

// BuildQueryGORM build datapaging for GORM DB instance
func (pagination *Datapaging) BuildQueryGORM(db *gorm.DB) *gorm.DB {

	if pagination.WithLimit() {
		db = db.Limit(pagination.Limit)
	}

	if pagination.WithPageOffset() {
		db = db.Offset(pagination.GetOffset())
	}

	if pagination.WithOrderBy() {
		db = db.Order(pagination.OrderBy[0])
	}

	if pagination.WithDateBetween() {
		db = db.Where("created_at > ? AND created_at < ?", pagination.DateEarliest.Unix(),
			pagination.DateLatest.Unix())
	}

	return db
}

// BuildQuery will add the pagination syntax into the raw sqlQuery
func (pagination *Datapaging) BuildQuery(sqlQuery string) string {

	if pagination.WithOrderBy() {
		sqlQuery = sqlQuery + " ORDER BY"
		for i, order := range pagination.OrderBy {
			if len(pagination.OrderBy)-i == 1 {
				sqlQuery = sqlQuery + " " + order
			} else {
				sqlQuery = sqlQuery + " " + order + ","
			}
		}
	}

	if pagination.WithLimit() {
		sqlQuery = sqlQuery + " LIMIT " + strconv.Itoa(pagination.Limit)
	}

	if pagination.WithPageOffset() {
		pageOffset := (pagination.Limit * (pagination.Page)) - pagination.Limit
		sqlQuery = sqlQuery + " OFFSET " + strconv.Itoa(pageOffset)
	}

	return sqlQuery
}
