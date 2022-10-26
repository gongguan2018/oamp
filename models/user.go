package models

type Username struct {
	ID        int    `gorm:"column:id" json:"id"`
	Username  string `gorm:"column:username" json:"username"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	Password  string `gorm:"column:password" json:"-"`
	Useremail string `gorm:"column:user_email" json:"useremail"`
	Userrole  string `gorm:"column:user_role" json:"userrole"`
	State     int    `gorm:"column:state" json:"state"`
}

func (*Username) TableName() string {
	return "oamp_user_info"
}

//判断登录用户名是否存在
func ExistsUser(u string) (bool, error) {
	var userinfo Username
	err := db.Select("username").Where("username = ?", u).Find(&userinfo).Error
	if err != nil {
		return false, err
	}
	if userinfo.Username != "" {
		return true, nil
	}
	return false, nil
}
func ExistsUserID(id int) (bool, error) {
	var userinfo Username
	err := db.Select("id").Where("id = ?", id).Find(&userinfo).Error
	if err != nil {
		return false, err
	}
	if userinfo.ID > 0 {
		return true, nil
	}
	return false, nil
}

//添加用户
func AddUser(data map[string]interface{}) error {
	err := db.Model(&Username{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

//修改用户
func EditUser(id int, data map[string]interface{}) error {
	var userinfo Username
	err := db.Model(&userinfo).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func GetUsername() ([]*Username, error) {
	var userinfo []*Username
	err := db.Select("id", "username", "nickname", "user_email", "user_role", "state").Find(&userinfo).Error
	if err != nil {
		return nil, err
	}
	return userinfo, nil
}

//删除用户
func DeleteUser(id int) error {
	var userinfo Username
	err := db.Where("id = ?", id).Delete(&userinfo).Error
	if err != nil {
		return err
	}
	db.Raw("SET @i=0;").Scan(&userinfo)
	db.Raw("UPDATE oamp_user_info SET id=(@i:=@i+1);").Scan(&userinfo)
	db.Raw("ALTER TABLE oamp_user_info auto_increment=1;").Scan(&userinfo)
	return nil

}
