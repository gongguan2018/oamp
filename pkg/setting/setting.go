package setting

import (
	"oamp/global"
	"time"

	"gopkg.in/ini.v1"
)

type Server struct {
	HttpPort     int
	ReadTimeout  time.Duration //默认是纳秒单位
	WriteTimeout time.Duration
}

//声明变量,值为实例化后的空结构体
var ServerSetting = &Server{}

type database struct {
	User        string
	Password    string
	Host        string
	DbName      string
	Port        string
	TablePrefix string
}

var DatabaseSetting = &database{}

type app struct {
	RunMode   string
	JwtSecret string
	PageSize  int
}

var AppSetting = &app{}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &redis{}

type log struct {
	Level         string
	Format        string
	Directory     string
	ShowLine      bool
	StacktraceKey string
	LogToConsole  bool
	EncodeLevel   string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
}

var LogSetting = &log{}

//编写函数读取配置文件
func Setup() {
	//读取配置文件
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		global.Log.Error("err")
	}
	//通过MapTo将配置文件log映射到结构体LogSetting中
	if err := cfg.Section("log").MapTo(LogSetting); err != nil {
		global.Log.Error(err.Error())
	}
	//将配置文件通过MapTo映射为结构体字段
	err = cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		global.Log.Error("err")
	}
	//将时间设置为秒数,配置中为60,因此为60s
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	//映射配置文件到结构体
	err = cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		global.Log.Error("err")
	}
	//映射配置文件到结构体
	err = cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		global.Log.Error("err")
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	if err := cfg.Section("redis").MapTo(RedisSetting); err != nil {
		global.Log.Error("err")
	}
}
