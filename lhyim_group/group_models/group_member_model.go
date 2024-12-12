package group_models

import "lhyim_server/common/models"

type GroupMemberModel struct {
	models.Model
	GroupID         uint       `json:"groupID"`                       //群ID
	GroupModel      GroupModel `gorm:"ForeignKey:GroupID" json:"-"`   //群
	UserID          uint       `json:"userID"`                        //用户ID
	Role            int8       `json:"role"`                          //角色 1群主 2管理员 3普通成员
	MemberNickname  string     `gorm:"size:32" json:"memberNickname"` //群内昵称
	ProhibitionTime *int       `json:"prohibitionTime"`               //禁言时间 单位分钟

}
