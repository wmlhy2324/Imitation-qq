package chat_models

import "lhyim_server/common/models"

type TopUserModel struct {
	models.Model
	UserID    uint `json:"UserID"`
	TopUserID uint `json:"TopUserID"`
}
