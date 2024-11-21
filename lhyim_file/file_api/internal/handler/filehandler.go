package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"lhyim_server/common/response"
	"lhyim_server/lhyim_file/file_api/internal/logic"
	"lhyim_server/lhyim_file/file_api/internal/svc"
	"lhyim_server/lhyim_file/file_api/internal/types"
	"lhyim_server/lhyim_file/file_model"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		file, fileHead, err := r.FormFile("file")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		//文件大小限制
		//size := float64(fileHead.Size) / float64(1024) / float64(1024)
		//if size > svcCtx.Config.FileSize {
		//	response.Response(r, w, nil, errors.New("图片超过限制大小,最大为200kb"))
		//	return
		//}
		//文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".")
		var suffix string //后缀
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}
		if utils.List(svcCtx.Config.BlackList, suffix) {
			response.Response(r, w, nil, errors.New("文件非法,请上传图片"))
			return
		}
		fileData, _ := io.ReadAll(file)
		fileHash := utils.MD5(fileData)
		l := logic.NewFileLogic(r.Context(), svcCtx)
		resp, err := l.File(&req)
		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", fileHash).Error
		if err == nil {
			resp.Src = fileModel.WebPath()
			logx.Infof("重复文件")
			response.Response(r, w, resp, err)
			return
		}
		//文件重名
		//拿用户信息
		userResponse, err := svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{UserIdList: []uint32{uint32(req.UserID)}})
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		dirName := fmt.Sprintf("%d_%s", req.UserID, userResponse.UserInfo[uint32(req.UserID)].NickName)
		dirPath := path.Join(svcCtx.Config.UploadDir, "file", dirName)

		_, err = os.ReadDir(dirPath)
		if err != nil {
			logx.Error(err)
			os.MkdirAll(dirPath, 0666)
		}
		NewfileModel := file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileHead.Filename,
			Size:     fileHead.Size,
			Hash:     fileHash,
			Uid:      uuid.New(),
		}
		filePath := path.Join(dirPath, fmt.Sprintf("%s.%s", NewfileModel.Uid, suffix))

		NewfileModel.Path = filePath
		//filename := fileHead.Filename

		err = os.WriteFile(filePath, fileData, 0666)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return

		}
		err = svcCtx.DB.Create(&NewfileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		resp.Src = NewfileModel.WebPath()
		response.Response(r, w, resp, err)

	}
}
