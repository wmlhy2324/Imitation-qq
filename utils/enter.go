package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
)

func List(list []string, key string) (ok bool) {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false

}
func ListByRegex(list []string, key string) (ok bool) {
	for _, s := range list {
		regex, err := regexp.Compile(s)
		if err != nil {
			logx.Error(err)
			return false
		}
		if regex.MatchString(key) {
			return true
		}
	}

	return false
}
func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}
func GetFilePrefix(filename string) (prefix string) {
	nameList := strings.Split(filename, ".")
	for i := 0; i < len(nameList)-1; i++ {
		if i == len(nameList)-2 {
			prefix += nameList[i]
		} else {
			prefix += nameList[i] + "."
		}
	}
	return
}

// 去重
func DeduplicatoionList[T string | int | uint | uint32](req []T) (response []T) {
	i32map := make(map[T]bool)
	for _, v := range req {
		if !i32map[v] {
			i32map[v] = true
		}
	}
	for key, _ := range i32map {
		response = append(response, key)
	}
	return response
}
