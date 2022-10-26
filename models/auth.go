package models

import (
	"errors"
	"oamp/pkg/util"
)

type UserInfo struct {
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func ExistsUserName(username string) (bool, error) {
	var user UserInfo
	err := db.Select("username").Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, errors.New("sql语句错误或者用户名没有发现")
	}
	if user.Username == "" {
		return false, errors.New("用户名不存在")
	}
	return true, nil
}

//下面是通过mysql自带的decode进行解密
//func CheckPassword(password string) (bool, error) {
//	var user UserInfo
/*
   通过Raw执行原生sql,通过select语句中的decode进行password解密,后面的
   as是将字段名指定为别名password，为什么要定义别名呢?
   因为结构体中的字段名和数据库字段名一一对应，只有查询到匹配的字段才可以映射到结构体中去,
   如果不指定别名,那么查询到的字段就是decode((select password from oamp_user_info WHERE id = 1),'5182086abcD#')
   结构体中无此字段，因为也不能映射过去,会出现扫描错误
   Scan()需要配合Raw一起用，功能跟Find()差不多,如需深入研究可参考源码
*/
//	db.Raw("select decode((select password from oamp_user_info WHERE id = 1),'5182086abcD#') as password;").Scan(&user)
//	if user.Password != password {
//		return false, errors.New("用户密码输入错误")
//	}
//	return true, nil

//}
func CheckPassword(username, password string) (bool, error) {
	var user UserInfo
	err := db.Select("password").Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, err
	}
	if util.AesDecrypt(user.Password) != password {
		return false, errors.New("密码不正确")
	}
	return true, nil
}
