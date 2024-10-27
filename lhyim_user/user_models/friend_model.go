package user_models

import (
	"gorm.io/gorm"
	"lhyim_server/common/models"
)

type FriendModel struct {
	models.Model
	SendUserID     uint      `json:"sendUserID"` //发起验证方
	SendUserModel  UserModel `gorm:"ForeignKey:SendUserID" json:"-"`
	RecvUserID     uint      `json:"recvUserID"` //接收验证方
	RecvUserModel  UserModel `gorm:"ForeignKey:RecvUserID" json:"-"`
	SendUserNotice string    `gorm:"size:128" json:"sendUserNotice"` //发起给接收方的备注
	RecvUserNotice string    `gorm:"size:128" json:"recvUserNotice"` //接收给验证方的备注
}

func (f *FriendModel) IsFriend(db *gorm.DB, A, B uint) bool {
	err := db.Take(&f, "(send_user_id = ? and recv_user_id = ?) or (recv_user_id = ? and send_user_id = ?)", A, B, A, B).Error
	if err == nil {
		return true
	}
	return false
}
func (f *FriendModel) Friends(db *gorm.DB, UserID uint) (list []FriendModel) {
	db.Find(&f, "send_user_id = ? or recv_user_id = ? ", UserID)

	return
}
func (f *FriendModel) GetUserNotice(UserID uint) string {
	if UserID == f.SendUserID {
		return f.SendUserNotice
	}
	if UserID == f.RecvUserID {
		return f.RecvUserNotice
	}
	return ""
}
