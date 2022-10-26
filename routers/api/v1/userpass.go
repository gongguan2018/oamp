package v1

import (
	"encoding/json"
	"net/http"
	"oamp/global"
	"oamp/models"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type Userpass struct {
	Username  string `json:"username" binding:"required,max=30,ne=0"`
	Password  string `json:"password" binding:"required,max=30,ne=0"`
	IPAddress string `json:"ipaddress" binding:"required,ip"`
}

//添加用户名和密码
func AddUserPass(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     []Userpass
	)
	//将gin接收的参数与结构体进行绑定
	err := c.Bind(&form)
	if err != nil {
		global.Log.Error(err.Error())
		response.Response(http.StatusBadRequest, errcode.INVALID_PARAMS, nil)
		return
	}
	//如果form的长度为0,说明切片没有元素,说明请求接口的时候没有传递任何内容
	if len(form) == 0 {
		response.Response(http.StatusBadRequest, errcode.INVALID_PARAMS_USER_PASS_NOT_NIL, nil)
		return
	}
	//将form进行遍历,将密码字段的值进行加密,采用对称加密方式
	for k, v := range form {
		EncryptCode := util.AesEncrypt(v.Password) //加密
		form[k].Password = EncryptCode             //加密后重新赋值给切片中的密码字段
	}
	//将结构体切片form序列化为字节切片
	byteSlice, err := json.Marshal(&form)
	if err := models.UserPass(byteSlice); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}

//获取用户名和密码信息
func GetUserPass(c *gin.Context) {
	var response = app.Gin{C: c}
	//实例化结构体字段,c.Param获取url的链接参数信息
	userpass := &Userpass{IPAddress: com.StrTo(c.Param("ipaddress")).String()}
	err, data := models.GetUsernamePass(userpass.IPAddress)
	if err != nil {
		response.Response(http.StatusBadRequest, errcode.ERROR_IP_NOT_EXISTS, nil)
		return
	}
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}

//删除用户名和密码
func DeleteUserPass(c *gin.Context) {
	var response = app.Gin{C: c}
	userpass := &Userpass{IPAddress: com.StrTo(c.Param("ipaddress")).String()}
	err := models.DeleteUserPass(userpass.IPAddress)
	if err != nil {
		response.Response(http.StatusBadRequest, errcode.ERROR_IP_NOT_EXISTS, nil)
		return
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}
