package v1

import (
	"errors"
	"fmt"
	"net/http"
	"oamp/global"
	"oamp/models"
	"oamp/pkg/app"
	"oamp/pkg/errcode"
	"oamp/pkg/gredis"
	"oamp/pkg/setting"
	"oamp/pkg/util"
	"oamp/service/system_service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type AddSystemForm struct {
	SystemType string    `form:"system_type" valid:"Required;MaxSize(20)"`
	IPAddress  string    `form:"ip_address"  valid:"Required;MaxSize(20)"`
	Hostname   string    `form:"hostname" valid:"Required;MaxSize(20)"`
	Created    string    `form:"created" valid:"Required;MaxSize(20)"`
	Remarks    string    `form:"remarks" valid:"Required;MinSize(1)"`
	CreatedAt  time.Time `gorm:"createdat"`
	UpdatedAt  time.Time `gorm:"updatedat"`
}

//添加系统信息
func AddSystem(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     AddSystemForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return //return表示结束当前函数，并返回指定值
	}
	//字段验证通过，现在判断表中是否存在相同数据
	addService := system_service.SystemInfo{
		SystemType: form.SystemType,
		IPAddress:  form.IPAddress,
		Hostname:   form.Hostname,
		Created:    form.Created,
		Remarks:    form.Remarks,
	}
	//判断数据库是否存在相同IP
	existsIP, err := addService.ExistsIP()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if existsIP {
		response.Response(http.StatusBadRequest, errcode.IP_EXISTS, nil)
		return
	}
	err = addService.AddSystem()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	//添加后获取所有记录条数
	count, err := models.GetSystemTotal()
	if err != nil {
		global.Log.Error(err.Error())
	}
	//将全部记录进行运算后,删除最后添加的记录所在页的Redis的key
	if err := delRedisKey(count); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}

type GetSystemForm struct {
	ID      int `form:"id" valid:"Min(0)"`
	PageNum int `form:"pagenum" valid:"Min(-1)"`
}

//获取系统信息
func GetSystem(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     = GetSystemForm{ID: com.StrTo(c.Param("id")).MustInt(), PageNum: util.GetPage(c)}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return //return表示结束当前函数，并返回指定值
	}
	//如果form.ID不为0,说明获取的是指定id的系统信息,否则获取的是全部系统信息
	if form.ID != 0 {
		getService := system_service.SystemInfo{
			ID: form.ID,
		}
		//判断是否存在此ID
		existsID, err := getService.ExistsID()
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		if !existsID {
			response.Response(http.StatusBadRequest, errcode.ID_NOT_EXISTS, nil)
			return
		}
		getData, err := getService.GetSystem()
		if err != nil {
			global.Log.Error(err.Error())
			response.Response(http.StatusBadRequest, errcode.ERROR_GET_DATA_FAIL, nil)
			return
		}
		response.Response(http.StatusOK, errcode.SUCCESS, getData)
	} else {
		//实例化结构体SystemInfo,PageNum表示当前页
		getService := system_service.SystemInfo{PageNum: form.PageNum}
		getData, err := getService.GetSystems()
		if err != nil {
			global.Log.Error(err.Error())
			response.Response(http.StatusBadRequest, errcode.ERROR_GET_DATA_FAIL, nil)
			return
		}
		//获取总数
		getSystemTotal, err := getService.GetSystemTotal()
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		data := make(map[string]interface{})
		data["List"] = getData
		data["Total"] = getSystemTotal
		data["PageSize"] = setting.AppSetting.PageSize
		response.Response(http.StatusOK, errcode.SUCCESS, data)
	}
}

type EditSystemForm struct {
	ID         int       `form:"id" valid:"Required;Min(1)"`
	SystemType string    `form:"system_type" valid:"Required;MaxSize(20)"`
	IPAddress  string    `form:"ip_address"  valid:"Required;MaxSize(20)"`
	Hostname   string    `form:"hostname" valid:"Required;MaxSize(20)"`
	Created    string    `form:"created" valid:"Required;MaxSize(20)"`
	Remarks    string    `form:"remarks" valid:"Required;MinSize(1)"`
	CreatedAt  time.Time `gorm:"createdat"`
	UpdatedAt  time.Time `gorm:"updatedat"`
}

//编辑系统信息
func EditSystem(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     = EditSystemForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	//判断是否存在此条记录，根据ID进行判断
	editService := system_service.SystemInfo{
		ID:         form.ID,
		SystemType: form.SystemType,
		IPAddress:  form.IPAddress,
		Hostname:   form.Hostname,
		Created:    form.Created,
		Remarks:    form.Remarks,
	}
	existsID, err := editService.ExistsID()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if !existsID {
		response.Response(http.StatusBadRequest, errcode.ID_NOT_EXISTS, nil)
		return
	}
	//如果存在此ID,调用函数编辑
	err = editService.EditSystem()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	//编辑完成后需要删除此ID对应的redis的key,这样在点击页面的时候才能获取最新信息
	delIDKey := int64(form.ID)
	if err := delRedisKey(delIDKey); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}

type DeleteSystemForm struct {
	ID int `form:"id" valid:"Required;Min(1)"`
}

//删除系统信息
func DeleteSystem(c *gin.Context) {
	var (
		response = app.Gin{C: c}
		form     = DeleteSystemForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != errcode.SUCCESS {
		response.Response(httpCode, errCode, nil)
		return
	}
	deleteService := system_service.SystemInfo{
		ID: form.ID,
	}
	existsID, err := deleteService.ExistsID()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	if !existsID {
		response.Response(http.StatusBadRequest, errcode.ID_NOT_EXISTS, nil)
		return
	}
	err = deleteService.DeleteSystem()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	delIDKey := int64(form.ID)
	if err := delRedisKey(delIDKey); err != nil {
		global.Log.Error(err.Error())
	}
	response.Response(http.StatusOK, errcode.SUCCESS, nil)
}

/*
  在编辑和添加信息后,redis中存在缓存无法实时更新,因此需要将redis的key删除
  因此需要根据点击时候的页数来删除对应的key,这样节约性能,无需每页重新读取数据库,只请求对应页的即可
  通过strconv.ParseFloat将字符串转换为float,为什么除以10?根据当前ID或者总的记录数除以10,取得的整数就表示当前的页数
  例如:如果所有的总记录数为85,每页显示的数据量为10,因此85实际对应的就是第9页,因为数据查询(Offset)的时候0对应的是第一页,因此实际上85/10后的8即表示
  第9页,但是又有一个问题,当页数为10的整数倍的时候,比如10、20、30,此时10/10得到的是1,1对应的是第二页,但是实际10是在第一页中,因为每页显示10条,因此此时
  删除的就是第二页的key,还是无法实时更新,因此当页数为10的整数倍的时候,需要将key调整为理论上的减去1才是正确的
*/
func delRedisKey(i int64) error {
	var splicing []string
	//获取浮点数,go中默认除法后取的是整数,不利于后期运算,因此转为浮点数,保留一位小数
	recordCount, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(i)/float64(10)), 64)
	//将获取的浮点数转换为字符串形式,保留一位小数
	fToStr := strconv.FormatFloat(recordCount, 'f', 1, 64)
	//截取字符串中.后面的字符,主要用来判断是否为0,从而知道页数是不是10的整数倍,如果是整数倍则为0
	splitStr := strings.Split(fToStr, ".")[1]
	//判断是否与0相等,如果不等则将当前页的值作为redis的key删除
	if splitStr != "0" {
		splicing = []string{"PAGENUM", strings.Split(fToStr, ".")[0]}
	} else {
		//如果与0相等,那么需要将当前页的值减去1作为redis的key删除
		splicing = []string{"PAGENUM", strconv.FormatInt(int64(recordCount)-1, 10)}
	}
	//拼接字符串作为要删除的key
	deleteKey := strings.Join(splicing, "_")
	delResult, err := gredis.Delete(deleteKey)
	if err != nil {
		return err
	}
	if !delResult {
		return errors.New("redis的key删除失败")
	}
	return nil
}
