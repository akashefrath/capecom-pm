package models

type Client struct {
	BaseModel

	Name    string
	Email   *string
	Phone   *string
	Address *string
}
