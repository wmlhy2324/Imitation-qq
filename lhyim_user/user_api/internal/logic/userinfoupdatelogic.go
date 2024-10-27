package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"
	"lhyim_server/lhyim_user/user_models"
	"lhyim_server/utils/maps"
)

type UserInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(req *types.UserInfoUpdateRequest) (resp *types.UserInfoUpdateResponse, err error) {
	// todo: add your logic here and delete this line
	userMaps := maps.RefToMap(*req, "user")
	if len(userMaps) != 0 {
		var user user_models.UserModel
		err = l.svcCtx.DB.Take(&user, req.UserID).Error
		if err != nil {
			return nil, errors.New("用户不存在")
		}
		err = l.svcCtx.DB.Model(&user).Updates(userMaps).Error
		if err != nil {
			logx.Error(userMaps)
			logx.Error(user)
			return nil, errors.New("用户信息更新失败")
		}
	}
	userConfMaps := maps.RefToMap(*req, "user_conf")
	if len(userConfMaps) != 0 {
		var userConf user_models.UserConfModel
		err = l.svcCtx.DB.Take(&userConf, "user_id = ?", req.UserID).Error
		if err != nil {
			return nil, errors.New("用户配置不存在")
		}

		verification, ok := userConfMaps["verification_question"]
		if ok {
			delete(userConfMaps, "verification_question")
			data := ctype.VerifcationQuestion{}

			maps.MapToStrcut(verification.(map[string]interface{}), &data)
			l.svcCtx.DB.Model(&userConf).Updates(&user_models.UserConfModel{
				VerificationQuestion: &data,
			})
		}
		err = l.svcCtx.DB.Model(&userConf).Updates(userConfMaps).Error
		if err != nil {
			logx.Error(userConfMaps)
			logx.Error(userConf)
			return nil, errors.New("用户配置信息更新失败")
		}
	}
	return
}
