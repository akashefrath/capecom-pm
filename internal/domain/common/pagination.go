package common

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
func (p *Pagination) BuildPaginationQuery() (string, []interface{}) {
	p.Normalize()
	// Return the query string AND the arguments separately
	return " LIMIT ? OFFSET ? ", []interface{}{p.Limit, p.Offset()}
}

func (p *Pagination) HasMore(total int64) bool {
	return int64(p.Page*p.Limit) < total
}
