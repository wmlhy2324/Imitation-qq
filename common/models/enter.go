package models

type Model struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
type PageInfo struct {
	Page  int    `form:"page,optional"`
	Limit int    `form:"limit,optional"`
	Sort  string `form:"sort,optional"`
	Key   string `form:"key,optional"`
}
