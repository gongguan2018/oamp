package main

import (
	"fmt"
	"net/http"
	"oamp/models"
	"oamp/pkg/gredis"
	"oamp/pkg/logging"
	"oamp/pkg/setting"
	"oamp/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	setting.Setup()
	logging.InitLog(setting.LogSetting.Level)
	//初始化连接数据库的函数,注意一定要在setting之后,因为setting.Setup需要加载配置文件
	models.Setup()
	gredis.Setup()
	//设置运行模式,默认情况下就是debug,生产环境要用release
	gin.SetMode(setting.AppSetting.RunMode)
	//调用routers包下的路由模块的InitRouter函数,返回*gin.Engine
	router := routers.InitRouter()
	r := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
		//每次左移1位，相当于乘以2,1<<20=1*2^20,也就是1048字节
	}
	//启动监听服务
	r.ListenAndServe()
}
