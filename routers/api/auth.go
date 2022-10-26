package api

import (
	"net/http"
	"oamp/global"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/util"
	"oamp/service/auth_service"

	"github.com/gin-gonic/gin"
)

//注意：字段名一定要大写,否则获取不到
type auth struct {
	Username string `form:"username" valid:"Required;MaxSize(20)"`
	Password string `form:"password" valid:"Required;MaxSize(50)"`
}

//用户名和密码的请求格式为json
func GetAuth(c *gin.Context) {
	var (
		/*
			            如果实例化Gin结构体的时候,不传递c过去，那么在执行c.JSON的时候就会提示指针异常
						因为Gin中的变量C的类型为*gin.Context，默认值为nil,通过nil去调用JSON肯定异常
		*/
		response = app.Gin{C: c}
		form     auth
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//实例化Auth结构体
	authService := auth_service.Auth{
		Username: form.Username,
		Password: form.Password,
	}
	//判断是否存在用户名
	exists, err := authService.ExistsUserName()
	if err != nil {
		global.Log.Error(err.Error())
	}
	//如果exists为false
	if !exists {
		response.Response(http.StatusBadRequest, errcode.ERR_USER_NOT_EXISTS, nil)
		return
	}
	//判断是否存在密码
	existsPass, err := authService.ExistsPassword()
	if err != nil {
		global.Log.Error(err.Error())
	}
	if !existsPass {
		response.Response(http.StatusBadRequest, errcode.ERROR_PASS, nil)
		return
	}
	//创建token
	token, err := util.CreateToken(form.Username, form.Password)
	if err != nil {
		global.Log.Error(err.Error())
	}
	data := make(map[string]interface{})
	data["token"] = token
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}
