package logic

import (
	"context"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserValidListLogic {
	return &UserValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserValidListLogic) UserValidList(req *types.FriendValidRequest) (resp *types.FriendValidResponse, err error) {
	// todo: add your logic here and delete this line
	fvs, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendVerifyModel{},
		list_query.Option{
			PageInfo: models.PageInfo{
				Page:  req.Page,
				Limit: req.Limit,
			},
			Where:   l.svcCtx.DB.Where("send_user_id = ? or recv_user_id = ?", req.UserID, req.UserID),
			Preload: []string{"RecvUserModel.UserConfModel", "SendUserModel.UserConfModel"},
		})
	var list []types.FriendValidInfo
	for _, v := range fvs {
		info := types.FriendValidInfo{
			AddtionalMessage: v.AdditionalMessage,
			ID:               v.ID,
			CreateAt:         v.CreatedAt.String(),
		}
		if v.SendUserID == req.UserID {
			//我是发起方
			info.UserID = v.RecvUserID
			info.Nickname = v.RecvUserModel.Nickname
			info.Avatar = v.RecvUserModel.Avatar
			info.Verification = v.RecvUserModel.UserConfModel.Verification
			info.Status = v.SendStatus
			info.Flag = "send"
		}
		if v.RecvUserID == req.UserID {
			//我是接收方
			info.UserID = v.SendUserID
			info.Nickname = v.SendUserModel.Nickname
			info.Avatar = v.SendUserModel.Avatar
			info.Verification = v.SendUserModel.UserConfModel.Verification
			info.Status = v.RecvStatus
			info.Flag = "rev"
		}

		if v.VerificationQuestion != nil {
			info.VerficationQuestion = &types.VerificationQuestion{
				Problem1: v.VerificationQuestion.Problem1,
				Problem2: v.VerificationQuestion.Problem2,
				Problem3: v.VerificationQuestion.Problem3,
				Answer1:  v.VerificationQuestion.Answer1,
				Answer2:  v.VerificationQuestion.Answer2,
				Answer3:  v.VerificationQuestion.Answer3,
			}
		}
		list = append(list, info)

	}
	return &types.FriendValidResponse{
		List:  list,
		Count: count,
	}, nil
}
