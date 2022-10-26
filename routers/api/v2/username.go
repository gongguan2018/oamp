package v2

import (
	"net/http"
	"oamp/global"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/util"
	"oamp/service/user_service"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type Username struct {
	ID        int    `form:"id"`
	Username  string `form:"username" valid:"Required;MaxSize(20)"`
	Password  string `form:"password" valid:"Required;MaxSize(20)"`
	Useremail string `form:"useremail" valid:"Required;MaxSize(50)"`
	Userrole  string `form:"userrole" valid:"Required;MaxSize(100)"`
	Nickname  string `form:"nickname" valid:"Required;MaxSize(20)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

//添加用户名和密码
func AddUsername(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     Username
	)
	//将gin的参数与结构体进行绑定
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//检查用户名是否存在,如果存在则拒绝添加
	userService := user_service.Username{
		Username:  form.Username,
		Password:  util.AesEncrypt(form.Password),
		Useremail: form.Useremail,
		Userrole:  form.Userrole,
		Nickname:  form.Nickname,
		State:     form.State,
	}
	existsUser, err := userService.ExistsUser()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if existsUser {
		response.Response(http.StatusBadRequest, errcode.USER_EXISTS, nil)
		return
	}
	if err := userService.AddUser(); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}

type Edituser struct {
	ID        int    `form:"id" valid:"Required;Min(1)"`
	Username  string `form:"username" valid:"MaxSize(20)"`
	Password  string `form:"password" valid:"MaxSize(20)"`
	Useremail string `form:"useremail" valid:"MaxSize(50)"`
	Userrole  string `form:"userrole" valid:"MaxSize(100)"`
	Nickname  string `form:"nickname" valid:"MaxSize(20)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

//修改用户名和密码
func EditUsername(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     = Edituser{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	//将gin的参数与结构体进行绑定
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//检查用户名是否存在,如果存在则拒绝添加
	userService := user_service.Username{
		ID:        form.ID,
		Username:  form.Username,
		Password:  util.AesEncrypt(form.Password),
		Useremail: form.Useremail,
		Userrole:  form.Userrole,
		Nickname:  form.Nickname,
		State:     form.State,
	}
	existsID, err := userService.ExistsID()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if !existsID {
		response.Response(http.StatusBadRequest, errcode.ID_NOT_EXISTS, nil)
		return
	}
	if err := userService.EditUser(); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)

}

//删除用户名和密码
func DeleteUsername(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     = Username{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	userService := user_service.Username{
		ID: form.ID,
	}
	existsID, err := userService.ExistsID()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if !existsID {
		response.Response(http.StatusBadRequest, errcode.ID_NOT_EXISTS, nil)
		return
	}
	if err := userService.DeleteUser(); err != nil {
		response.Response(http.StatusBadRequest, errcode.DEFAULT_USER_CANNOT_BE_DELETE, nil)
		global.Log.Error(err.Error())
		return
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)

}
func GetUsername(c *gin.Context) {
	var response = app.Gin{C: c}
	getUsername := user_service.Username{}
	data, err := getUsername.GetUsername()
	if err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}
