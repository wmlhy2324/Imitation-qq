package group_models

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"lhyim_server/common/models"
	"time"
)

type GroupMemberModel struct {
	models.Model
	GroupID         uint       `json:"groupID"`                       //群ID
	GroupModel      GroupModel `gorm:"ForeignKey:GroupID" json:"-"`   //群
	UserID          uint       `json:"userID"`                        //用户ID
	Role            int8       `json:"role"`                          //角色 1群主 2管理员 3普通成员
	MemberNickname  string     `gorm:"size:32" json:"memberNickname"` //群内昵称
	ProhibitionTime *int       `json:"prohibitionTime"`               //禁言时间 单位分钟

}

func (gm *GroupMemberModel) GetProhibitionTime(client *redis.Client, db *gorm.DB) *int {
	if gm.ProhibitionTime == nil {
		return nil
	}
	t, err := client.TTL(fmt.Sprintf("prohibition_%d", gm.ID)).Result()
	if err != nil || t == -2*time.Second {
		//查不到说明过期了
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}

	res := int(t.Seconds())
	return &res
}
