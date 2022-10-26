package models

import (
	"fmt"
	"oamp/global"
	"oamp/pkg/setting"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//声明全局变量db
var db *gorm.DB

func Setup() {
	var (
		err                                       error
		dbName, user, password, host, tableprefix string
	)
	dbName = setting.DatabaseSetting.DbName
	user = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tableprefix = setting.DatabaseSetting.TablePrefix
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,        //设置表名为单数形式
			TablePrefix:   tableprefix, //设置表前缀
		},
	})
	if err != nil {
		global.Log.Error(err.Error())
	}
	//使用数据库连接池
	sqlDB, err := db.DB()
	//设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	//设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	//设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
}
