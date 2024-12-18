package logic

import (
	"context"
	"errors"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/utils/maps"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUpdateLogic) GroupUpdate(req *types.GroupUpdataRequest) (resp *types.GroupUpdateResponse, err error) {
	//只能是群主或者管理员才能调用
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&groupMember, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("群不存在或者用户不是群成员")
	}
	if !(groupMember.Role == 1 || groupMember.Role == 2) {
		return nil, errors.New("群信息只能是群主或者管理员更新")
	}
	groupMaps := maps.RefToMap(*req, "conf")
	if len(groupMaps) != 0 {
		verificationQuestion, ok := groupMaps["verification_question"]
		if ok {
			delete(groupMaps, "verification_question")
			data := ctype.VerifcationQuestion{}
			maps.MapToStrcut(verificationQuestion.(map[string]any), &data)
			l.svcCtx.DB.Model(&groupMember.GroupModel).Updates(&group_models.GroupModel{
				VerificationQuestion: &data,
			})
		}
		err = l.svcCtx.DB.Model(&groupMember.GroupModel).Updates(groupMaps).Error
		if err != nil {
			logx.Error(err)
			return nil, errors.New("群信息更新失败")
		}
	}
	return
}
