package app

import (
	"net/http"
	"oamp/global"
	"oamp/pkg/errcode"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

//将gin的参数绑定到结构体中，并进行参数检查
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	//Bind()、ShouldBind()将查询参数,http Head,数据格式(json,xml)绑定到结构体中
	//Bind()和ShouldBind()区别:ShouldBind没有绑定成功不报错，就是空值,Bind会报错
	err := c.Bind(form)
	if err != nil {
		global.Log.Error(err.Error())
		return http.StatusBadRequest, errcode.ERROR_PARAMS_BIND_FAIL
	}
	valid := validation.Validation{}
	//检验form结构体字段中设定的valid是否有效，返回bool和err
	check, err := valid.Valid(form)
	if err != nil {
		global.Log.Error(err.Error())
		return http.StatusInternalServerError, errcode.INVALID_PARAMS
	}
	//如果无效,也就是不符合要求
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, errcode.INVALID_PARAMS
	}
	return http.StatusOK, errcode.SUCCESS
}
