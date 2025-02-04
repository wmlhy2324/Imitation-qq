// Code generated by goctl. DO NOT EDIT.
package types

type AddGroupRequest struct {
	UserID               uint                  `header:"User-ID"`
	GroupID              uint                  `json:"groupId"`
	Verify               string                `json:"verify,optional"`      //验证消息
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"` //问题和答案
}

type AddGroupResponse struct {
}

type GroupMemberInfo struct {
	UserID         uint   `json:"userId"`
	UserNickName   string `json:"userNickName"`
	Avatar         string `json:"avatar"`
	IsOnline       bool   `json:"isOnline"`
	Role           int8   `json:"role"`
	MemberNickname string `json:"memberNickname"`
	CreateAt       string `json:"createAt"`
	NewMsgDate     string `json:"newMsgDate"`
}

type GroupMyListResponse struct {
	List  []GroupMyResponse `json:"list"`
	Count int               `json:"count"`
}

type GroupMyResponse struct {
	GroupID          uint   `json:"groupId"`
	GroupAvatar      string `json:"groupAvatar"`
	GroupTitle       string `json:"groupTitle"`
	GroupMemberCount int    `json:"groupMemberCount"`
	Role             int8   `json:"role"`
	Mole             int8   `json:"mole"`
}

type GroupSearch struct {
	GroupID         uint   `json:"groupId"`
	Title           string `json:"title"`
	Abstract        string `json:"abstract"`
	Avatar          string `json:"avatar"`
	IsInGroup       bool   `json:"isInGroup"` //我是否在群里
	UserCount       int    `json:"userCount"` //总数
	UserOnlineCount int    `json:"userOnlineCount"`
}

type GroupSessionList struct {
	GroupID       uint   `json:"groupId"`
	Title         string `json:"title"`
	Avatar        string `json:"avatar"`
	NewMsgDate    string `json:"newMsgDate"`    //最新消息
	NewMsgPreview string `json:"newMsgPreview"` //最新的消息内容
	IsTop         bool   `json:"isTop"`         //是否置顶
}

type GroupValidRequest struct {
	UserID  uint `header:"User-ID"`
	GroupID uint `path:"id"`
}

type GroupValidResponse struct {
	Verification         int8                 `json:"verification"`         //好友验证
	VerificationQuestion VerificationQuestion `json:"verificationQuestion"` //问题和答案
}

type GroupValidiInfo struct {
	ID                   uint                  `json:"id"` //验证id
	GroupID              uint                  `json:"groupId"`
	UserID               uint                  `json:"userId"`
	UserNickname         string                `json:"userNickname"`
	UserAvatar           string                `json:"userAvatar"`
	Status               int8                  `json:"status"` //状态
	AddtionalMessage     string                `json:"addtionalMessage"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"`
	Title                string                `json:"title"` //群名
	CreateAt             string                `json:"createAt"`
	Type                 int8                  `json:"type"` //1 加 2退
}

type GroupfriendsList struct {
	UserId    uint   `json:"userId"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	IsInGroup bool   `json:"isInGroup"` //是否在群里面
}

type UserInfo struct {
	UserID   uint   `header:"User-ID"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1,optional" conf:"problem1"`
	Problem2 *string `json:"problem2,optional" conf:"problem2"`
	Problem3 *string `json:"problem3,optional" conf:"problem3"`
	Answer1  *string `json:"answer1,optional" conf:"answer1"`
	Answer2  *string `json:"answer2,optional" conf:"answer2"`
	Answer3  *string `json:"answer3,optional" conf:"answer3"`
}

type GroupChatRequest struct {
	UserID uint `header:"User-ID"`
}

type GroupChatResponse struct {
}

type GroupCreateRequest struct {
	UserID     uint   `header:"User-ID"`
	Mode       int8   `json:"mode,optional"` //1直接创建 2选人创建
	Name       string `json:"name,optional"`
	IsSearch   bool   `json:"isSearch,optional"`   //是否可以搜到
	Size       int    `json:"size,optional"`       //群规模
	UserIDList []uint `json:"userIdList,optional"` //用户id列表
}

type GroupCreateResponse struct {
}

type GroupHistoryDeleteRequest struct {
	UserID    uint   `header:"User-ID"`
	GroupID   uint   `path:"id"`
	Page      int    `form:"page,optional"`
	Limit     int    `form:"limit,optional"`
	MsgIDList []uint `json:"msgIdList"`
}

type GroupHistoryDeleteResponse struct {
}

type GroupHistoryRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `path:"id"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type GroupHistoryResponse struct {
}

type GroupInfoRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `path:"id"`
}

type GroupInfoResponse struct {
	GroupID           uint       `json:"groupId"`
	Title             string     `json:"title"`
	Abstract          string     `json:"abstract"`
	MemberCount       int        `json:"memberCount"`
	MemberOnlineCount int        `json:"memberOnlineCount"`
	Avatar            string     `json:"avatar"`
	Creator           UserInfo   `json:"creator"`         //群主
	AdminList         []UserInfo `json:"adminList"`       //管理员列表
	Role              int8       `json:"role"`            //1群主，2群管理员，3群成员
	IsProhibition     bool       `json:"isProhibition"`   //群禁言
	ProhibitionTime   *int       `json:"prohibitionTime"` //自己的禁言时间
}

type GroupMemberAddRequest struct {
	UserID       uint   `header:"User-ID"`
	ID           uint   `json:"id"` //群id
	MemberIDList []uint `json:"memberIdlist"`
}

type GroupMemberAddResponse struct {
}

type GroupMemberNicknameUpdateRequest struct {
	UserID   uint   `header:"User-ID"`
	ID       uint   `json:"id"` //群id
	MemberID uint   `json:"memberId"`
	Nickname string `json:"nickname"`
}

type GroupMemberNicknameUpdateResponse struct {
}

type GroupMemberRemoveRequest struct {
	UserID   uint `header:"User-ID"`
	ID       uint `form:"id"` //群id
	MemberID uint `form:"memberId"`
}

type GroupMemberRemoveResponse struct {
}

type GroupMemberRequest struct {
	UserID uint   `header:"User-ID"`
	ID     uint   `form:"id"`
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
	Sort   string `form:"sort,optional"`
}

type GroupMemberResponse struct {
	List  []GroupMemberInfo `json:"list"`
	Count int               `json:"count"`
}

type GroupMemberRoleUpdateRequest struct {
	UserID   uint `header:"User-ID"`
	ID       uint `json:"id"` //群id
	MemberID uint `json:"memberId"`
	Role     int8 `json:"role"`
}

type GroupMemberRoleUpdateResponse struct {
}

type GroupMyRequest struct {
	UserID uint `header:"User-ID"`
	Mode   int8 `form:"mode"` //1我创建的群聊 2我加入的群聊
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type GroupProhibitionRequest struct {
	UserID          uint `header:"User-ID"`
	MemberID        uint `json:"memberId"`
	ProhibitionTime *int `json:"prohibitionTime,optional"` //分钟
	GroupID         uint `json:"groupId"`
}

type GroupProhibitionResponse struct {
}

type GroupRemoveRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `path:"id"`
}

type GroupRemoveResponse struct {
}

type GroupSearchListResponse struct {
	List  []GroupSearch `json:"list"`
	Count int           `json:"count"`
}

type GroupSearchRequest struct {
	UserID uint   `header:"User-ID"`
	Key    string `form:"key"` //群id和昵称
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
}

type GroupSessionRequest struct {
	UserID uint `header:"User-ID"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type GroupSessionResponse struct {
	List  []GroupSessionList `json:"list"`
	Count int                `json:"count"`
}

type GroupTopRequest struct {
	UserID  uint `header:"User-ID"`
	GroupID uint `json:"groupId"` //需要置顶的群id
	IsTop   bool `json:"isTop"`   //true置顶 false取消置顶
}

type GroupTopResponse struct {
}

type GroupUpdataRequest struct {
	UserID               uint                  `header:"User-ID"`
	ID                   uint                  `json:"id"`
	IsSearch             *bool                 `json:"isSearch,optional" conf:"is_search"`
	Verification         *int8                 `json:"verification,optional" conf:"verification"`
	IsInvite             *bool                 `json:"isInvite,optional" conf:"is_invite"`
	IsTemporarySession   *bool                 `json:"isTemporarySession,optional" conf:"is_temporary_session"` //是否开启临时会话
	IsProhibition        *bool                 `json:"isProhibition,optional" conf:"is_prohibition"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional" conf:"verification_question"`
	Avatar               *string               `json:"avatar,optional" conf:"avatar"`
	Abstract             *string               `json:"abstrat,optional" conf:"abstract"`
	Title                *string               `json:"title,optional" conf:"title"`
}

type GroupUpdateResponse struct {
}

type GroupValidListRequest struct {
	UserID uint `header:"User-ID"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type GroupValidListResponse struct {
	List  []GroupValidiInfo `json:"list"`
	Count int               `json:"count"`
}

type GroupValidStatusRequest struct {
	UserID  uint `header:"User-ID"`
	ValidID uint `json:"validId"` //验证id
	Status  int8 `json:"status"`  //状态id
}

type GroupValidStatusResponse struct {
}

type GroupfriendsListRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `form:"id"` //群id
}

type GroupfriendsListResponse struct {
	List  []GroupfriendsList `json:"list"`
	Count int                `json:"count"`
}
