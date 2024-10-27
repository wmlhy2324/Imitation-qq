package etcd

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"lhyim_server/core"
	"strings"
)

func DeliveryAddress(etcdaddr string, serviceName string, addr string) {
	list := strings.Split(addr, ":")
	if len(list) != 2 {
		logx.Errorf("addr format error ", addr)
		return
	}
	if list[0] == "0.0.0.0" {
		ip := netx.InternalIp()

		addr = strings.ReplaceAll(addr, "0.0.0.0", ip)
	}
	client := core.InitEtcd(etcdaddr)
	_, err := client.Put(context.Background(), serviceName, addr)
	if err != nil {
		logx.Errorf("put etcd error %s", err.Error())
		return
	}
	logx.Infof("put etcd success %s ,%s", serviceName, addr)
}
func GetServiceAddr(etcdaddr string, serviceName string) string {
	client := core.InitEtcd(etcdaddr)
	resp, err := client.Get(context.Background(), serviceName)
	if err == nil && len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value)
	}
	return ""
}
