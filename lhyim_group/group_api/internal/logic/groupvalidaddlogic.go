package logic

import (
	"context"
	"errors"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"
	"lhyim_server/lhyim_group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidAddLogic {
	return &GroupValidAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidAddLogic) GroupValidAdd(req *types.AddGroupRequest) (resp *types.AddGroupResponse, err error) {
	//加群
	//判断是否在群
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err == nil {
		return nil, errors.New("不需要重复加群")
	}
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, req.GroupID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}
	resp = new(types.AddGroupResponse)
	var verifyModel = group_models.GroupVerifyModel{
		AdditionalMessage: req.Verify,
		GroupID:           req.GroupID,
		UserID:            req.UserID,
		Status:            0,
		Type:              1,
	}
	switch group.Verification {
	case 0: //不允许添加
		return nil, errors.New("不允许任何人加群")
	case 1: //允许任何人
		verifyModel.Status = 1
		var groupMember = group_models.GroupMemberModel{
			GroupID: req.GroupID,
			UserID:  req.UserID,
			Role:    3,
		}
		l.svcCtx.DB.Create(&groupMember)
	case 2: //需要验证问题

	case 3: //需要回答问题
		if req.VerificationQuestion != nil {
			verifyModel.VerificationQuestion = &ctype.VerifcationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
				Answer1:  req.VerificationQuestion.Answer1,
				Answer2:  req.VerificationQuestion.Answer2,
				Answer3:  req.VerificationQuestion.Answer3,
			}
		}
	case 4:
		var count int //记录回答问题的个数
		if req.VerificationQuestion != nil && group.VerificationQuestion != nil {
			if group.VerificationQuestion.Answer1 != nil && req.VerificationQuestion.Answer1 != nil {
				if *group.VerificationQuestion.Answer1 == *req.VerificationQuestion.Answer1 {
					count++
				}
			}
			if group.VerificationQuestion.Answer2 != nil && req.VerificationQuestion.Answer2 != nil {
				if *group.VerificationQuestion.Answer2 == *req.VerificationQuestion.Answer2 {

					count++
				}
			}
			if group.VerificationQuestion.Answer3 != nil && req.VerificationQuestion.Answer3 != nil {
				if *group.VerificationQuestion.Answer3 == *req.VerificationQuestion.Answer3 {
					count++
				}
			}

			if count != group.ProblemCount() {
				return nil, errors.New("答案错误")
			}
			//直接加群
			verifyModel.Status = 1
			verifyModel.VerificationQuestion = group.VerificationQuestion
			//把用户加到群里面
			var groupMember = group_models.GroupMemberModel{
				GroupID: req.GroupID,
				UserID:  req.UserID,
				Role:    3,
			}
			l.svcCtx.DB.Create(&groupMember)

		} else {
			return nil, errors.New("答案错误")
		}

	default:

	}
	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		return nil, err
	}
	return
}
