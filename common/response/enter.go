package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Body struct {
	Code uint32      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		r := &Body{
			Code: 0,
			Data: resp,
			Msg:  "成功",
		}
		httpx.WriteJson(w, http.StatusOK, r)

		return
	}
	CodeMsg := uint32(7)
	//ErrorMsg := "服务器错误"
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: CodeMsg,
		Data: nil,
		Msg:  err.Error(),
	})
}
