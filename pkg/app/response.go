package app

import (
	"oamp/pkg/errcode"

	"github.com/gin-gonic/gin"
)

//声明一个名称为Gin的结构体
type Gin struct {
	//结构体类型为gin的指针类型,注意C要大写，因为其他文件要调用
	C *gin.Context
}

//httpCode: http状态码
//errCode:  错误代码
//JSON()是gin.Context的方法,因此需要通过g.C调用gin.Context进而调用JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  errcode.GetMsg(errCode),
		"data": data,
	})
	return
}
