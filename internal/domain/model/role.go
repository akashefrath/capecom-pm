package models

type Role struct {
	BaseModel
	Name string `db:"name" json:"name"`
}
