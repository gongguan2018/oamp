package jwt

import (
	"net/http"
	"oamp/global"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

//此中间件用于jwt和gin接口对接
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			code     int
			response = app.Gin{C: c}
		)
		code = errcode.SUCCESS
		//接收来自请求头的key的信息，key的名字为Authorization，可自定义,需要与请求的时候传递的名称一致
		token := c.Request.Header.Get("Authorization")
		//		token := c.Query("token")
		if token == "" {
			code = errcode.INVALID_PARAMS
			global.Log.Error("token不存在,请携带正确的token进行访问")
		} else {
			//解析token
			claims, err := util.ParseToken(token)
			if err != nil {
				global.Log.Error(err.Error())
				code = errcode.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				global.Log.Error("Token已过期,请重新生成")
				code = errcode.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != errcode.SUCCESS {
			response.Response(http.StatusUnauthorized, code, nil)
			//Abor()表示中断,如果上述条件都不符合要求，将中断请求
			c.Abort()
			return
		}
		//Next()表示挂起，先执行其余函数，最后执行c.Next()后面的
		//如果不设置Abort()那么上面的报错信息出现后，还是会继续执行Next()后面的
		c.Next()
		/*
			          例如：在路由文件router.go中，配置了获取用户信息的路由/user/info，并且设置了认证中间件jwt.JWT()
					  如果在jwt.go中设置了Abort()，那么当执行出错的时候，比如没有输入token，此时就会中断执行，不会继续调用
					  后面的执行程序api.UserInfo,如果没有设置Abort()，那么在执行出错后，打印了错误信息后，通过c.Next()还会执行后面的处理
					  程序api.UserInfo，相当于出现了两个结果,jwt就失去了原有的效力
		*/
	}
}
