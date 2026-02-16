package response

type APIResponse struct {
	Success bool        `json:"success" default:"true"`
	Message string      `json:"message,omitempty"`
	Func    string      `json:"func,omitempty"`
	Entity  string      `json:"entity,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PageMeta struct {
	Page    int     `json:"page,omitempty"`
	Limit   int     `json:"limit,omitempty"`
	Total   int64   `json:"total,omitempty"`
	Pages   int     `json:"pages,omitempty"`
	NextID  *string `json:"next_id,omitempty"`
	HasNext bool    `json:"has_next,omitempty"`
}
