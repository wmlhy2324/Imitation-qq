package user_models

import (
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
)

type UserConfModel struct {
	models.Model
	UserID               uint                       `json:"userID"`
	UserModel            UserModel                  `gorm:"ForeignKey:UserID" json:"-"`
	RecallMessage        *string                    `gorm:"size:32" json:"recallMessage"` //撤回消息
	FriendOnline         bool                       `json:"friendOnline"`                 //好友上线提醒
	Sound                bool                       `json:"sound"`                        //声音提醒
	SecureLink           bool                       `json:"secureLink"`                   //安全链接
	SavePwd              bool                       `json:"savePwd"`                      //保存密码
	SearchUser           int8                       `json:"searchUser"`                   //搜索用户 0不能被搜索 1 用过id搜索 2 通过用户名找到我
	Verification         int8                       `json:"Verification"`                 //好友验证 0不允许任何人添加 1允许任何人添加 2需要验证 3需要验证且回答问题 4需要验证且正确回答问题
	VerificationQuestion *ctype.VerifcationQuestion `json:"verificationQuestion"`         //验证问题
	Online               bool                       `json:"online"`                       //在线状态
}

func (uc UserConfModel) ProblemCount() (c int) {
	if uc.VerificationQuestion != nil {
		if uc.VerificationQuestion.Problem1 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem2 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem3 != nil {
			c += 1
		}
	}
	return c
}
