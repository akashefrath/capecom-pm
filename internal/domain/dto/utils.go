package dto

type UtilOption struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type UtilsResponse struct {
	Roles        []UtilOption `json:"roles"`
	UserGroups   []UtilOption `json:"user_groups"`
	Designations []UtilOption `json:"designations"`
	Departments  []UtilOption `json:"departments"`
	Clients      []UtilOption `json:"clients"`
	TicketTypes  []UtilOption `json:"ticket_types"`
	Users        []UtilOption `json:"users"`
}
