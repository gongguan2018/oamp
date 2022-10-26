package models

import (
	"encoding/json"
	"errors"
	"oamp/pkg/gredis"
	"oamp/pkg/util"
	"oamp/service/cache_service"
	"time"
)

type UserpassInfo struct {
	ID        int       `gorm:"column:id" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	IPAddress string    `gorm:"column:ipaddress" json:"ipaddress"`
	CreatedAt time.Time `gorm:"column:createdat" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat" json:"updatedat"`
}

//函数接收传递来的IP地址和结构体切片
func getUserPass(ip string, userpass []UserpassInfo) error {
	var user []UserpassInfo
	//根据IP地址查询数据,根据模型映射到结构体切片user中
	db.Where("ipaddress = ?", ip).Find(&user)
	//如果切片的元素数量大于0，说明查询到了此IP数据
	if len(user) > 0 {
		/*
			将从数据库查询后的切片user的长度与传递过来的切片userpass长度做对比
			如果user < userpass,说明在新增用户密码，将len(user)作为起始索引，将切片截取，新的切片元素即为新增的元素
			如果user > userpass,说明在删除用户密码，将len(userpass)作为起始索引，截取切片，新的切片即为要删除的元素内容
			如果user = userpass,说明数据在修改，没有新增和删除,通过for循环,变量范围为user切片的长度,根据数据库中的ipaddress和id字段，进行循环更新数据
		*/
		if len(user) < len(userpass) {
			newSlice := userpass[len(user):]
			err := db.Create(&newSlice).Error //可根据结构体切片创建数据
			if err != nil {
				return err
			}
			b, e := DeleteKey(ip)
			if e != nil {
				return err
			}
			if b {
				return nil
			} else {
				return errors.New("key删除失败")
			}
		} else if len(user) > len(userpass) {
			newSlice := user[len(userpass):]
			err := db.Delete(&newSlice).Error //可根据结构体切片删除数据
			if err != nil {
				return err
			}
			db.Raw("SET @i=0;").Scan(&UserpassInfo{})
			db.Raw("UPDATE oamp_userpass_info SET id=(@i:=@i+1);").Scan(&UserpassInfo{})
			db.Raw("ALTER TABLE oamp_userpass_info auto_increment=1;").Scan(&UserpassInfo{})
			b, e := DeleteKey(ip)
			if e != nil {
				return err
			}
			if b {
				return nil
			} else {
				return errors.New("key删除失败")
			}
		} else {
			for i := 0; i < len(user); i++ {
				err := db.Model(UserpassInfo{}).Where("ipaddress = ? AND id = ?", ip, user[i].ID).Updates(userpass[i]).Error
				if err != nil {
					return err
				}
			}
			b, e := DeleteKey(ip)
			if e != nil {
				return e
			}
			if b {
				return nil
			} else {
				return errors.New("key删除失败")
			}
		}
	} else {
		err := db.Create(&userpass).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 UserPass函数接收参数为字节切片，通过Unmarshal将此字节切片反序列化到结构体切片中
*/
func UserPass(b []byte) error {
	var userpass []UserpassInfo
	err := json.Unmarshal(b, &userpass)
	if err != nil {
		return err
	}
	//获取反序列化后的结构体切片中的第一个元素的IP地址，因为切片中的元素IP地址都是一样，因此获取第一个第二个都一样
	ip := userpass[0].IPAddress
	//调用函数，将结构体切片和IP地址传递到函数中
	if err := getUserPass(ip, userpass); err != nil {
		return err
	}
	//db.Debug().Create(&userpass)
	return nil
}

//查询用户名和密码信息，返回类型为错误和结构体切片
func GetUsernamePass(ip string) (error, []UserpassInfo) {
	var userpass []UserpassInfo
	IPAddressKey := cache_service.GetIPKey(ip)
	if gredis.Exists(IPAddressKey) {
		data, err := gredis.Get(IPAddressKey)
		if err != nil {
			return nil, nil
		} else {
			json.Unmarshal(data, &userpass)
			return nil, userpass
		}
	}
	err := db.Select("username", "password", "ipaddress").Where("ipaddress = ?", ip).Find(&userpass).Error
	if err != nil {
		return err, nil
	}
	//如果切片长度等于0,说明没有查询到了数据
	if len(userpass) == 0 {
		return errors.New("没有查询到IP数据"), nil
	}
	//从数据库查询的password字段为加密后的,因此需要通过AesDecrypt解密
	for k, v := range userpass {
		DecryptCode := util.AesDecrypt(v.Password) //解密
		userpass[k].Password = DecryptCode         //解密后重新赋值
	}
	gredis.Set(IPAddressKey, userpass, 3600)
	return nil, userpass

}

//删除用户名和密码
func DeleteUserPass(ip string) error {
	var userpass UserpassInfo
	err := db.Select("ipaddress").Where("ipaddress = ?", ip).Find(&userpass).Error
	if err != nil {
		return err
	}
	if userpass.IPAddress == "" {
		return errors.New("没有此条IP记录")
	}
	if err := db.Where("ipaddress = ?", ip).Delete(&userpass).Error; err != nil {
		return err
	}
	db.Raw("SET @i=0;").Scan(&UserpassInfo{})
	db.Raw("UPDATE oamp_userpass_info SET id=(@i:=@i+1);").Scan(&UserpassInfo{})
	db.Raw("ALTER TABLE oamp_userpass_info auto_increment=1;").Scan(&UserpassInfo{})
	b, e := DeleteKey(ip)
	if e != nil {
		return err
	}
	if b {
		return nil
	} else {
		return errors.New("key删除失败")
	}
	return nil
}
func DeleteKey(ip string) (bool, error) {
	IPAddressKey := cache_service.GetIPKey(ip)
	if gredis.Exists(IPAddressKey) {
		//调用函数Delet删除redis的key
		deleteResult, err := gredis.Delete(IPAddressKey)
		if err != nil {
			return false, err
		}
		//如果为false，表示删除失败
		if !deleteResult {
			return false, nil
		}
		return true, nil
	}
	return true, nil
}

//func AddIPAddress(data map[string]interface{}) error {
//	err := db.Create(&UserpassInfo{
//		Username:  data["username"].(string),
//		Password:  data["password"].(string),
//		IPAddress: data["ipaddress"].(string),
//	}).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
