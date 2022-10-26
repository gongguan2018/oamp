package app

import (
	"oamp/global"

	"github.com/astaxie/beego/validation"
)

//通过validation执行校验结构体字段
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		global.Log.Error(err.Error())
	}
	return
}
