package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"lhyim_server/common/etcd"

	"github.com/zeromicro/go-zero/core/conf"

	"net/http"
	"regexp"
	"strings"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func failResponse(w http.ResponseWriter, msg string) {

	date, _ := json.Marshal(Response{Code: 7, Msg: msg})
	w.Write(date)
}
func auth(authUrl string, w http.ResponseWriter, r *http.Request) (ok bool) {

	authReq, err := http.NewRequest("POST", authUrl, nil)
	if err != nil {

		failResponse(w, "新请求生成失败")
		return false

	}

	authReq.Header = r.Header
	authReq.Header.Set("ValidPath", r.URL.Path)

	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {

		failResponse(w, "执行新请求失败")
		return false
	}
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data *struct {
			UserID uint `json:"userId"`
			Role   int  `json:"role"`
		} `json:"data"`
	}
	var authResponse Response

	byteDate, err := io.ReadAll(authRes.Body)

	if err != nil {

		failResponse(w, "读取新请求body失败")
		return false
	}

	authErr := json.Unmarshal(byteDate, &authResponse)
	fmt.Println("测试authResponse:", authResponse)
	if authErr != nil {
		fmt.Println(authErr)
		failResponse(w, "解析新请求body失败")
		return false
	}

	//认证失败

	if authResponse.Code != 0 {
		w.Write(byteDate)
		return false
	}
	if authResponse.Data != nil {
		r.Header.Set("User-ID", fmt.Sprintf("%d", authResponse.Data.UserID))
		r.Header.Set("Role", fmt.Sprintf("%d", authResponse.Data.Role))
	}
	return true

}
func proxy(proxyaddr string, w http.ResponseWriter, r *http.Request) {

	byteDate, _ := io.ReadAll(r.Body)

	proxyReq, err := http.NewRequest(r.Method, proxyaddr, bytes.NewBuffer(byteDate))

	proxyReq.Header = r.Header
	proxyReq.Header.Del("ValidPath")
	if err != nil {
		fmt.Println(err)
		failResponse(w, "服务异常")
		return

	}

	remoteAddr := strings.Split(r.RemoteAddr, ":")
	proxyReq.Header.Set("X-Forwarded-For", remoteAddr[0])
	response, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		fmt.Println(err)
		failResponse(w, "服务异常")
		return

	}
	io.Copy(w, response.Body)

}
func gateway(w http.ResponseWriter, r *http.Request) {
	//匹配请求前缀
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	addrList := regex.FindStringSubmatch(r.URL.Path)

	if len(addrList) != 2 {
		failResponse(w, "err path not found")
		return
	}
	service := addrList[1]
	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	if addr == "" {
		fmt.Println("服务未找到", service)
		failResponse(w, "err service not found")
		return
	}
	remoteAddr := strings.Split(r.RemoteAddr, ":")

	//请求认证服务地址
	authaddr := etcd.GetServiceAddr(config.Etcd, "auth_api")
	if authaddr == "" {

		failResponse(w, "err auth not found")
		return

	}
	//获取请求地址

	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authaddr)
	ok := auth(authUrl, w, r)
	if !ok {
		return
	}

	//认证到此结束
	url := fmt.Sprintf("http://%s%s", addr, r.URL.String())
	proxy(url, w, r)
	logx.Infof("请求地址为%s 代理地址为%s", remoteAddr[0], url)

}

var configFile = flag.String("f", "settings.yaml", "the config file")

type Config struct {
	Addr string `json:"addr"`
	Etcd string `json:"etcd"`
	Log  logx.LogConf
}

var config Config

func main() {
	flag.Parse()
	logx.SetUp(config.Log)
	conf.MustLoad(*configFile, &config)
	http.HandleFunc("/", gateway)
	fmt.Println("gateway server start at ", config.Addr)
	http.ListenAndServe(config.Addr, nil)
}
