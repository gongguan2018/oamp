package user_service

import (
	"encoding/json"
	"errors"
	"oamp/models"
	"oamp/pkg/gredis"
	"oamp/service/cache_service"
)

type Username struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Useremail string `json:"useremail"`
	Userrole  string `json:"userrole"`
	State     int    `json:state`
}

//检查用户是否存在的方法
//注意: 方法作用于结构体,因此可直接通过u.Username访问
func (u *Username) ExistsUser() (bool, error) {
	return models.ExistsUser(u.Username)
}
func (u *Username) AddUser() error {
	data := getData(u)
	keys := cache_service.GetUserKeys()
	if gredis.Exists(keys) {
		_, err := gredis.Delete(keys)
		if err != nil {
			return err
		}
	}
	return models.AddUser(data)
}

//修改用户
func (u *Username) EditUser() error {
	data := getData(u)
	keys := cache_service.GetUserKeys()
	if gredis.Exists(keys) {
		_, err := gredis.Delete(keys)
		if err != nil {
			return err
		}
	}
	return models.EditUser(u.ID, data)
}
func (u *Username) ExistsID() (bool, error) {
	return models.ExistsUserID(u.ID)
}

//删除用户
func (u *Username) DeleteUser() error {
	if u.ID == 1 {
		return errors.New("此用户为默认用户,不可删除")
	}
	keys := cache_service.GetUserKeys()
	if gredis.Exists(keys) {
		_, err := gredis.Delete(keys)
		if err != nil {
			return err
		}
	}
	return models.DeleteUser(u.ID)
}

//获取用户信息
func (u *Username) GetUsername() ([]*models.Username, error) {
	var username []*models.Username
	keys := cache_service.GetUserKeys()
	if gredis.Exists(keys) {
		data, err := gredis.Get(keys)
		if err != nil {
			return nil, err
		} else {
			json.Unmarshal(data, &username)
			return username, nil
		}

	}
	users, err := models.GetUsername()
	if err != nil {
		return nil, err
	}
	gredis.Set(keys, users, 3600)
	return users, nil
}

/*
    通过if判断结构体字段是否为空,如果不为空,说明有数据,此时将数据添加进map中
	此功能主要用来分组更新,如果都有数据,说明更新全部,如果部分有数据,那么就更新有数据的部分字段
*/
func getData(u *Username) map[string]interface{} {
	if u.State != 1 {
		u.State = 0
	}
	data := make(map[string]interface{})
	if u.Username != "" {
		data["username"] = u.Username
	}
	if u.Password != "" {
		data["password"] = u.Password
	}
	if u.Useremail != "" {
		data["user_email"] = u.Useremail
	}
	if u.Userrole != "" {
		data["user_role"] = u.Userrole
	}
	if u.Nickname != "" {
		data["nickname"] = u.Nickname
	}
	data["state"] = u.State
	return data
}
