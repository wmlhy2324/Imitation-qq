package handler

import (
	"errors"
	"fmt"
	"io"
	"lhyim_server/common/response"
	"lhyim_server/lhyim_file/file_api/internal/logic"
	"lhyim_server/lhyim_file/file_api/internal/svc"
	"lhyim_server/lhyim_file/file_api/internal/types"
	"lhyim_server/utils"
	"lhyim_server/utils/random"
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
		//文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".")
		var suffix string //后缀
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}
		if !utils.List(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("文件非法,请上传图片"))
			return
		}
		//文件重名
		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)

		dir, err := os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		filePath := path.Join(svcCtx.Config.UploadDir, imageType, fileHead.Filename)
		imageData, _ := io.ReadAll(file)
		//filename := fileHead.Filename
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)
		resp.Url = "/" + filePath
		if InDir(dir, fileHead.Filename) {
			byteDate, _ := os.ReadFile(filePath)
			oldfilehash := utils.MD5(byteDate)
			newfilehash := utils.MD5(imageData)
			if oldfilehash == newfilehash {
				fmt.Println("两个文件一样")
				response.Response(r, w, resp, nil)
				return
			}
			//两个文件不一样
			//改名操作
			var prefix = utils.GetFilePrefix(fileHead.Filename)
			newPath := fmt.Sprintf("%s_%s.%s", prefix, random.RandStr(4), suffix)
			filePath = path.Join(svcCtx.Config.UploadDir, imageType, newPath)
			//如果改了名字还是重名就需要递归
		}

		err = os.WriteFile(filePath, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return

		}
		resp.Url = "/" + filePath
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
