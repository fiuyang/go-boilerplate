package response

import "gorm.io/gorm"

type Meta struct {
	Limit     int   `json:"limit"`
	Page      int   `json:"page"`
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}

type Pagination struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Meta    Meta        `json:"meta"`   
}

func Scopes(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
