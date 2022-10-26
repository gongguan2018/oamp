package cache_service

import (
	"oamp/pkg/cache"
	"strconv"
	"strings"
)

type System struct {
	ID int
}
type IPAddress struct {
	IPAddress string
}

//根据ID返回对应的key
func (s *System) GetSystemKey() string {
	return cache.CACHE_IP + "_" + strconv.Itoa(s.ID) //strconv.Itoa将整型转为字符串
}
func (s *System) GetSystemKeys() string {
	keys := []string{
		cache.CACHE_ALL,
		"LIST",
	}
	return strings.Join(keys, "_")
}
func GetIPKey(ip string) string {
	return cache.CACHE_IPADDRESS + "_" + ip
}
func (s *System) GetPageKey(page int) string {
	return cache.CACHE_PAGENUM + "_" + strconv.Itoa(page)
}
