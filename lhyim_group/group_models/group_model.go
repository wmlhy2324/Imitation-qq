package group_models

import (
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
)

type GroupModel struct {
	models.Model
	Title                string                     `gorm:"size:32" json:"title"`     //群名称
	Avatar               string                     `gorm:"size:64" json:"avatar"`    //群头像
	Abstract             string                     `gorm:"size:256" json:"abstract"` //群简介
	Creator              uint                       `json:"creator"`                  //创建者
	Verification         int8                       `json:"Verification"`             //群验证 0不允许任何人添加 1允许任何人添加 2需要验证 3需要验证且回答问题 4需要验证且正确回答问题
	VerificationQuestion *ctype.VerifcationQuestion `json:"verificationQuestion"`     //验证问题
	IsSearch             bool                       `json:"isSearch"`                 //是否允许被搜索
	IsInvite             bool                       `json:"isInvite"`                 //是否允许被邀请
	IsTemporarySession   bool                       `json:"isTemporarySession"`       //是否是临时会话
	IsProhibition        bool                       `json:"isProhibition"`            //是否禁言
	Size                 int                        `json:"size"`                     //群成员数量
	MemberList           []GroupMemberModel         `gorm:"ForeignKey:GroupID" json:"memberList"`
}
