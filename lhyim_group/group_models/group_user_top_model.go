package group_models

import "lhyim_server/common/models"

type GroupUserTopModel struct {
	models.Model
	UserID  uint `json:"userID"`
	GroupID uint `json:"groupID"`
}
