package user_models

import (
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
)

type FriendVerifyModel struct {
	models.Model
	SendUserID           uint                       `json:"sendUserID"` //发起验证方
	SendUserModel        UserModel                  `gorm:"ForeignKey:SendUserID" json:"-"`
	RecvUserID           uint                       `json:"recvUserID"` //接收验证方
	RecvUserModel        UserModel                  `gorm:"ForeignKey:RecvUserID" json:"-"`
	Status               uint8                      `json:"status"`                            //状态 0未处理 1已同意 2已拒绝 3已忽略 4删除
	SendStatus           int8                       `json:"sendStatus"`                        //状态 0未处理 1已同意 2已拒绝 3已忽略 4删除
	RecvStatus           int8                       `json:"recvStatus"`                        //状态 0未处理 1已同意 2已拒绝 3已忽略 4删除
	AdditionalMessage    string                     `gorm:"size:128" json:"additionalMessage"` //附加消息
	VerificationQuestion *ctype.VerifcationQuestion `json:"verificationQuestion"`              //验证问题 为3和4的时候需要
}
