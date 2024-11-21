package handler

import (
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
	"lhyim_server/utils"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		file, fileHead, err := r.FormFile("image")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		imageType := r.FormValue("imageType")
		if imageType == "" {
			response.Response(r, w, nil, errors.New("imageType is empty"))
			return
		}
		switch imageType {
		case "avatar", "chat", "group_avatar":
		default:
			response.Response(r, w, nil, errors.New("imageType只能为avatar,chat,group_avatar"))
			return
		}
		//文件大小限制
		size := float64(fileHead.Size) / float64(1024) / float64(1024)
		if size > svcCtx.Config.FileSize {
			response.Response(r, w, nil, errors.New("图片超过限制大小,最大为200kb"))
			return
		}
		//文件上传用黑名单
		//文件后缀白名单

		nameList := strings.Split(fileHead.Filename, ".")
		var suffix string //后缀
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}
		if utils.List(svcCtx.Config.BlackList, suffix) {
			response.Response(r, w, nil, errors.New("文件非法,请上传正确文件格式"))
			return
		}
		//先去算hash
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)

		imageDate, _ := io.ReadAll(file)
		imageHash := utils.MD5(imageDate)
		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", imageHash).Error
		if err == nil {
			logx.Infof("重复文件")
			resp.Url = fileModel.WebPath()
			response.Response(r, w, resp, err)
			return
		}
		//拼接路径

		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)

		_, err = os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		NewfileModel := file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileHead.Filename,
			Size:     fileHead.Size,
			Hash:     utils.MD5(imageDate),
			Uid:      uuid.New(),
		}
		NewfileModel.Path = path.Join(dirPath, fmt.Sprintf("%s.%s", NewfileModel.Uid, suffix))

		//filename := fileHead.Filename

		err = os.WriteFile(NewfileModel.Path, imageDate, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return

		}
		//文件信息入库

		err = svcCtx.DB.Create(&NewfileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		resp.Url = NewfileModel.WebPath()

		response.Response(r, w, resp, err)

	}
}
func InDir(dir []os.DirEntry, file string) bool {
	for _, empty := range dir {
		if empty.Name() == file {
			return true
		}
	}
	return false
}
