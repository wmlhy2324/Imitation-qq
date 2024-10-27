package group_models

import (
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
)

type GroupVerifyModel struct {
	models.Model
	AdditionalMessage    string                     `gorm:"size:32" json:"additionalMessage"` //附加消息
	VerificationQuestion *ctype.VerifcationQuestion `json:"verificationQuestion"`             //验证问题 为3和4的时候需要
	GroupID              uint                       `json:"groupID"`                          //群ID
	GroupModel           GroupModel                 `gorm:"ForeignKey:GroupID" json:"-"`      //群

	UserID uint `json:"userID"` //用户ID
	Status int8 `json:"status"` //状态 0未处理 1已同意 2已拒绝 3已忽略
	Type   int8 `json:"type"`   //类型 1加群 2退群

}
