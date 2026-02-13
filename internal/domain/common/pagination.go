package common

import "fmt"

type Pagination struct {
	Page  int `json:"page" form:"page" query:"page"`
	Limit int `json:"limit" form:"limit" query:"limit"`
}

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

func (p *Pagination) Normalize() {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}

	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}

	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
func (p *Pagination) BuildPaginationQuery() string {
	p.Normalize()
	return fmt.Sprintf(" LIMIT %d OFFSET %d ", p.Limit, p.Offset())
}

func (p *Pagination) HasMore(total int64) bool {
	return int64(p.Page*p.Limit) < total
}
