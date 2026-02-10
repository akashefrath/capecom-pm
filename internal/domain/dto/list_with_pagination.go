package dto

type ListWithMeta struct {
	Data interface{} `json:"data"`
	Meta any         `json:"meta"`
}

type PaginationMeta struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"has_more"`
}
