package pagination

import (
	"math"

	"gorm.io/gorm"
)

// Param 分页参数
type Param struct {
	DB      *gorm.DB
	Page    int64
	Limit   int64
	OrderBy []string
	ShowSQL bool
}

// Paginator 分页返回
type Paginator struct {
	TotalRecord int64         `json:"total_record"`
	TotalPage   int64         `json:"total_page"`
	Offset      int64         `json:"offset"`
	Limit       int64         `json:"limit"`
	Page        int64         `json:"page"`
	PrevPage    int64         `json:"prev_page"`
	NextPage    int64         `json:"next_page"`
	Records     int64erface{} `json:"records"`
}

// Paging 分页
func Paging(p *Param, result int64erface{}) *Paginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int64
	var offset int64

	go countRecords(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	<-done

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int64(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, anyType int64erface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}
