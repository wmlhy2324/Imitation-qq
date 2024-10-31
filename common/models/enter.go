package models

import "time"

type Model struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type PageInfo struct {
	Page  int    `form:"page,optional"`
	Limit int    `form:"limit,optional"`
	Sort  string `form:"sort,optional"`
	Key   string `form:"key,optional"`
}
