package log_stash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Pusher struct {
	LogType    int8   `json:"logType"` //2操作日志 2运行日志
	IP         string `json:"ip"`
	UserID     uint   `json:"userID"`
	Level      string `json:"level"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Service    string `json:"service"` //记录微服务的名称
	client     *kq.Pusher
	items      []string
	ctx        context.Context
	isRequest  bool
	isHeader   bool
	isResponse bool
	request    string
	headers    string
	response   string
	count      int
}

func (p *Pusher) IsRequest() {
	p.isRequest = true
}
func (p *Pusher) IsHeader() {
	p.isHeader = true
}
func (p *Pusher) IsResponse() {
	p.isResponse = true
}
func (p *Pusher) GetResponse() bool {
	return p.isResponse
}
func (p *Pusher) SetRequest(r *http.Request) {
	// 请求头
	// 请求体
	// 请求路径，请求方法
	// 关于请求体的问题，拿了之后要还回去
	// 一定要在参数绑定之前调用

	method := r.Method
	path := r.URL.String()
	byteData, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(byteData))

	p.request = fmt.Sprintf(
		`<div class="log_request">
            <div class="log_request_head">
                <span class="log_request_method %s">%s</span>
                <span class="log_request_path">%s</span>
            </div>
            <div class="log_request_body">
                <pre class="log_json_body">%s</pre>
            </div>
        </div>`,
		strings.ToLower(method), method, path, string(byteData))
}
func (p *Pusher) SetHeader(r *http.Request) {
	byteData, _ := json.Marshal(r.Header)
	format := fmt.Sprintf(`<div class="log_request_header">
<div class="log_request_body">
<pre class="log_json_body">%s</pre>
</div>
</div>`, string(byteData))

	// 使用 fmt.Sprintf 进行格式化

	p.headers = format
}
func (p *Pusher) SetResponse(w string) {
	p.response = fmt.Sprintf("<div class=\"log_response\">\n <pre class=\"log_json_body\">%s</pre>\n</div>", w)
	p.Save(p.ctx)
}
func (p *Pusher) SetCtx(ctx context.Context) {
	p.ctx = ctx
}
func (p *Pusher) Info(title string) {
	p.Title = title
	p.Level = "info"
}
func (p *Pusher) Warning(title string) {
	p.Title = title
	p.Level = "warning"
}
func (p *Pusher) Err(title string) {
	p.Title = title
	p.Level = "err"
}
func (p *Pusher) Save(ctx context.Context) {
	if p.isResponse && p.count == 0 {
		p.count = 1
		return
	}
	if p.ctx == nil {
		p.ctx = ctx
	}
	if ctx == nil {
		p.ctx = context.Background()
	}
	if p.client == nil {
		return
	}
	var items []string
	if p.isRequest {
		items = append(items, p.request)
	}
	if p.isHeader {
		items = append(items, p.headers)
	}
	items = append(items, p.items...)
	if p.isResponse {
		items = append(items, p.response)
	}
	for _, item := range items {

		p.Content += item
	}
	p.items = []string{}

	userID := p.ctx.Value("userID")
	var ui int
	if userID != nil {
		ui, _ = strconv.Atoi(userID.(string))
	}
	clientIP, ok := p.ctx.Value("clientIP").(string)
	if ok {
		p.IP = clientIP
	}
	p.UserID = uint(ui)
	p.IP = clientIP
	if p.Level == "" {
		p.Level = "info"
	}
	byteData, err := json.Marshal(p)
	if err != nil {
		logx.Error(err)
		return
	}
	err = p.client.Push(ctx, string(byteData))
	if err != nil {
		logx.Error(err)
		return
	}
}
func (p *Pusher) SetItem(label string, val any) {

	p.setItem("info", label, val)
}
func (p *Pusher) SetItemInfo(label string, val any) {

	p.setItem("info", label, val)
}
func (p *Pusher) SetItemWarning(label string, val any) {

	p.setItem("warning", label, val)
}
func (p *Pusher) SetItemErr(label string, val any) {

	p.setItem("err", label, val)
}
func (p *Pusher) setItem(level string, label string, val any) {
	var str string
	switch value := val.(type) {
	case string:
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%s</div>", label, value)
	case int, int32, uint, uint8:
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%d</div>", label, value)
	default:
		byteData, _ := json.Marshal(val)
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%s</div>", label, string(byteData))
	}
	logItem := fmt.Sprintf("<div class=\"log_item %s\">%s</div>", level, str)
	p.items = append(p.items, logItem)

	if p.LogType == 3 {
		//如果是运行日志，那就调一下发送
		p.Save(p.ctx)
	}
}
func NewActionPusher(client *kq.Pusher, serviceName string) *Pusher {
	return NewPusher(client, serviceName, 2)
}
func NewRuntimePusher(client *kq.Pusher, serviceName string) *Pusher {
	return NewPusher(client, serviceName, 3)
}

func NewPusher(client *kq.Pusher, serviceName string, LogType int8) *Pusher {

	return &Pusher{
		Service: serviceName,
		client:  client,
		LogType: LogType,
	}

}
