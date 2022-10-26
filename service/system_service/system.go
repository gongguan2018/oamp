package system_service

import (
	"encoding/json"
	"oamp/global"
	"oamp/models"
	"oamp/pkg/gredis"
	"oamp/pkg/setting"
	"oamp/service/cache_service"
)

type SystemInfo struct {
	ID         int
	SystemType string
	IPAddress  string
	Hostname   string
	Created    string
	Remarks    string
	PageNum    int
	PageSize   int
}

func (s *SystemInfo) ExistsIP() (bool, error) {
	return models.ExistsIP(s.IPAddress)
}

//此函数用户在新增编辑删除的时候同步检查redis的key并删除key，实现数据实时更新
func CheckKeyDelete(id int) error {
	cache := cache_service.System{ID: id}
	//获取key
	key := cache.GetSystemKeys()
	//如果redis存在此key
	if gredis.Exists(key) {
		_, err := gredis.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SystemInfo) AddSystem() error {
	err := CheckKeyDelete(s.ID)
	if err != nil {
		global.Log.Error(err.Error())
	}
	data := make(map[string]interface{})
	data["system_type"] = s.SystemType
	data["ip_address"] = s.IPAddress
	data["hostname"] = s.Hostname
	data["created"] = s.Created
	data["remarks"] = s.Remarks
	return models.AddSystem(data)
}

//添加系统信息的同时，向表oamp_userpass_info中添加一条IP地址信息,避免在密码管理的时候查询不到而报错
//func (s *SystemInfo) AddIPAddress() error {
//	data := make(map[string]interface{})
//	data["username"] = ""
//	data["password"] = ""
//	data["ipaddress"] = s.IPAddress
//	return models.AddIPAddress(data)
//}
//判断是否存在系统id
func (s *SystemInfo) ExistsID() (bool, error) {
	return models.ExistsID(s.ID)
}

//编辑系统,根据id
func (s *SystemInfo) EditSystem() error {
	err := CheckKeyDelete(s.ID)
	if err != nil {
		global.Log.Error(err.Error())
	}
	data := make(map[string]interface{})
	data["system_type"] = s.SystemType
	data["ip_address"] = s.IPAddress
	data["hostname"] = s.Hostname
	data["created"] = s.Created
	data["remarks"] = s.Remarks
	return models.EditSystem(s.ID, data)
}

//删除系统,根据id
func (s *SystemInfo) DeleteSystem() error {
	err := CheckKeyDelete(s.ID)
	if err != nil {
		global.Log.Error(err.Error())
	}
	return models.DeleteSystem(s.ID)
}

//获取指定id的系统数据
func (s *SystemInfo) GetSystem() (*models.SystemInfo, error) {
	var cacheIP *models.SystemInfo
	cache := cache_service.System{ID: s.ID}
	//获取key
	key := cache.GetSystemKey()
	//如果redis存在此key
	if gredis.Exists(key) {
		//根据key获取value，value类型为[]byte
		data, err := gredis.Get(key)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, nil
		} else {
			//反序列化到结构体
			json.Unmarshal(data, &cacheIP)
			return cacheIP, nil
		}
	}
	//如果不存在key,直接去数据库中查找数据,ip为结构体的内存地址
	ip, err := models.GetSystem(s.ID)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, nil
	}
	//设置key
	gredis.Set(key, ip, 3600)
	return ip, nil
}

//获取全部系统信息
func (s *SystemInfo) GetSystems() ([]*models.SystemInfo, error) {
	cache := cache_service.System{}
	/*
		      如果s.PageNum小于0,说明请求的时候并没有指定页数,默认请求全部数据
			  PageNum的值通过gin的Query获取,在util的page.go中设定
	*/
	if s.PageNum == -1 {
		key := cache.GetSystemKeys()
		data, err := CheckExistsKey(key)
		if err != nil {
			return nil, err
		}
		if data == nil {
			ips, err := models.GetSystems()
			if err != nil {
				global.Log.Error(err.Error())
				return nil, nil
			}
			gredis.Set(key, ips, 3600)
			return ips, nil
		}
		return data, nil
	} else { //else表示s.PageNum的值不为-1,说明请求的时候携带了页数
		key := cache.GetPageKey(s.PageNum)
		data, err := CheckExistsKey(key)
		if err != nil {
			return nil, err
		}
		//如果不存在此key,那么data为nil
		if data == nil {
			ips, err := models.GetSystemsPage(s.PageNum, setting.AppSetting.PageSize)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, nil
			}
			gredis.Set(key, ips, 3600)
			return ips, nil
		}
		return data, nil
	}
	return nil, nil
}

//检查redis中是否存在此key
func CheckExistsKey(key string) ([]*models.SystemInfo, error) {
	var cacheIP []*models.SystemInfo
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			return nil, err
		} else {
			json.Unmarshal(data, &cacheIP)
			return cacheIP, nil
		}
	}
	return nil, nil
}
func (s *SystemInfo) GetSystemTotal() (int64, error) {
	return models.GetSystemTotal()
}
