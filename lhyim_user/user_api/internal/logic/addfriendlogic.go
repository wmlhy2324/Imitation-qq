package logic

import (
	"context"
	"errors"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendLogic) AddFriend(req *types.AddFriendRequest) (resp *types.AddFriendResponse, err error) {
	// todo: add your logic here and delete this line

	//如果是好久就不用加了
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, req.FriendID, req.UserID) {
		return nil, errors.New("你们已经是好友了")
	}
	var userConf user_models.UserConfModel
	err = l.svcCtx.DB.Take(&userConf, "user_id=?", req.FriendID).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	resp = new(types.AddFriendResponse)
	var verifyModel = user_models.FriendVerifyModel{
		SendUserID:        req.UserID,
		RecvUserID:        req.FriendID,
		AdditionalMessage: req.Verify,
	}
	//判断验证状态
	switch userConf.Verification {
	case 0: //不允许添加
		return nil, errors.New("该用户不允许任何人添加")
	case 1: //允许任何人
		verifyModel.RecvUserID = 1
		return nil, errors.New("已添加为好友")
	case 2: //需要验证问题

	case 3: //需要回答问题
		if req.VerificationQuestion != nil {
			verifyModel.VerificationQuestion = &ctype.VerifcationQuestion{
				Problem1: req.VerificationQuestion.Problem1,
				Problem2: req.VerificationQuestion.Problem2,
				Problem3: req.VerificationQuestion.Problem3,
				Answer1:  req.VerificationQuestion.Answer1,
				Answer2:  req.VerificationQuestion.Answer2,
				Answer3:  req.VerificationQuestion.Answer3,
			}
		}

	case 4: //需要回答问题
		//判断问题是否回答正确
		var count int //记录回答问题的个数
		if req.VerificationQuestion != nil && userConf.VerificationQuestion != nil {
			if userConf.VerificationQuestion.Answer1 != nil && req.VerificationQuestion.Answer1 != nil {
				if *userConf.VerificationQuestion.Answer1 == *req.VerificationQuestion.Answer1 {
					count++
				}
			}
			if userConf.VerificationQuestion.Answer2 != nil && req.VerificationQuestion.Answer2 != nil {
				if *userConf.VerificationQuestion.Answer2 == *req.VerificationQuestion.Answer2 {

					count++
				}
			}
			if userConf.VerificationQuestion.Answer3 != nil && req.VerificationQuestion.Answer3 != nil {
				if *userConf.VerificationQuestion.Answer3 == *req.VerificationQuestion.Answer3 {
					count++
				}
			}

			if count != userConf.ProblemCount() {
				return nil, errors.New("答案错误")
			}
			//直接加好友
			verifyModel.RecvUserID = 1
			verifyModel.VerificationQuestion = userConf.VerificationQuestion
			//加好友
			var userFriend = user_models.FriendModel{
				SendUserID: req.UserID,
				RecvUserID: req.FriendID,
			}
			l.svcCtx.DB.Create(&userFriend)
		}
	default:

	}
	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("添加好友失败")
	}
	return
}
