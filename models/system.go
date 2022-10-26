package models

import (
	"time"
)

type SystemInfo struct {
	ID         int       `gorm:"column:id" json:"id"`
	SystemType string    `gorm:"column:system_type" json:"system_type"`
	IPAddress  string    `gorm:"column:ip_address" json:"ip_address"`
	Hostname   string    `gorm:"column:hostname" json:"hostname"`
	Created    string    `gorm:"column:created" json:"created"`
	Remarks    string    `gorm:"column:remarks" json:"remarks"`
	CreatedAt  time.Time `gorm:"column:createdat" json:"createdat"`
	UpdatedAt  time.Time `gorm:"column:updatedat" json:"updatedat"`
}

func ExistsIP(ipaddress string) (bool, error) {
	var systeminfo SystemInfo
	err := db.Select("ip_address").Where("ip_address = ?", ipaddress).Find(&systeminfo).Error
	if err != nil {
		return false, err
	}
	if systeminfo.IPAddress == "" {
		return false, nil
	}
	return true, nil
}
func AddSystem(data map[string]interface{}) error {
	err := db.Create(&SystemInfo{
		SystemType: data["system_type"].(string), //类型断言
		IPAddress:  data["ip_address"].(string),
		Hostname:   data["hostname"].(string),
		Created:    data["created"].(string),
		Remarks:    data["remarks"].(string),
	}).Error
	if err != nil {
		return err
	}
	return nil
}
func ExistsID(id int) (bool, error) {
	var systeminfo SystemInfo
	err := db.Select("id").Where("id = ?", id).Find(&systeminfo).Error
	if err != nil {
		return false, err
	}
	if systeminfo.ID > 0 {
		return true, nil
	}
	return false, nil
}
func EditSystem(id int, data map[string]interface{}) error {
	var systeminfo SystemInfo
	err := db.Model(&systeminfo).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

/*
  根据数据库的id信息删除数据,db.Raw表示执行数据库原生SQL
  删除数据库id后,再新增数据会导致id出现断层,因此需要执行下面三条命令重新实现自增
  db.Raw("SET @i=0;").Scan(&SystemInfo{})
  db.Raw("UPDATE oamp_system_info SET id=(@i:=@i+1);").Scan(&SystemInfo{})
  db.Raw("ALTER TABLE oamp_system_info auto_increment=1;").Scan(&SystemInfo{})
*/
func DeleteSystem(id int) error {
	err := db.Where("id = ?", id).Delete(&SystemInfo{}).Error
	if err != nil {
		return err
	}
	db.Raw("SET @i=0;").Scan(&SystemInfo{})
	db.Raw("UPDATE oamp_system_info SET id=(@i:=@i+1);").Scan(&SystemInfo{})
	db.Raw("ALTER TABLE oamp_system_info auto_increment=1;").Scan(&SystemInfo{})
	return nil
}

//获取单条记录信息
func GetSystem(id int) (*SystemInfo, error) {
	var systeminfo SystemInfo
	err := db.Where("id = ?", id).Find(&systeminfo).Error
	if err != nil {
		return nil, err
	}
	return &systeminfo, nil
}

//获取全部记录信息
func GetSystems() ([]*SystemInfo, error) {
	var systeminfo []*SystemInfo
	err := db.Find(&systeminfo).Error
	if err != nil {
		return nil, err
	}
	return systeminfo, nil
}

//获取分页数据,pagenum表示传递的页码,pagesize表示显示的页的总数
/*
  pagenum * pagesize(每页显示记录数)什么意思呢？因为在执行分页查询的时候Offset表示要跳过的数据,limit表示要获取的数据条数,因此当页大小设置为10页
  的时候,第一条数据查询了10页数据,第二条数据需要跳过当前的10页,从第11页开始获取,以此类推,因此需要用当前页*每页的总数量(PageSize)
*/
func GetSystemsPage(pagenum, pagesize int) ([]*SystemInfo, error) {
	var systeminfo []*SystemInfo
	//Offset表示跳过的记录数,Limit表示获取的记录数
	//err := db.Offset(pagenum * pagesize).Limit(pagesize).Find(&systeminfo).Error
	//通过Offset和where的id都可以实现分页,Offset需要将表数据都遍历一次,性能有影响,因此建议使用where和id的形式来分页
	err := db.Where("id > ?", pagenum*pagesize).Limit(pagesize).Find(&systeminfo).Error
	if err != nil {
		return nil, err
	}
	return systeminfo, nil
}
func GetSystemTotal() (int64, error) {
	var (
		systeminfo SystemInfo
		count      int64
	)
	err := db.Model(&systeminfo).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
