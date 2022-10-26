package v2

import (
	"net/http"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/util"

	"github.com/gin-gonic/gin"
)

type encdec struct {
	Str string `form:"str" valid:"Required" json:"str"`
}

//加密
func Encryption(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     encdec
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//绑定成功后调用加密函数
	data := util.AesEncrypt(form.Str)
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}

//解密
func Decryption(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     encdec
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//绑定成功后调用解密函数
	data := util.AesDecrypt(form.Str)
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}
