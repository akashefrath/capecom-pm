package dto

type RoleResponse struct {
	ID   string `db:"uuid" json:"id"`
	Name string `db:"name" json:"name"`
	// Status string `json:"status"`
}
