package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

/*
  GetPage函数参数为*gin.Context,通过Query查询参数PageNum的值,如果传递的值大于0
  那么将此值减去1后返回,否则返回-1
  pageNum表示第几页,当传入的参数为1时候,表示查看第一页数据,但实际在进行db.OffSet运算的时候0才是第一页数据
  因此需要对当前传递的页数减去1
  如果不传递任何值或者传递的值等于0,返回-1,为什么返回-1？其实这个-1是自定义的,返回给/router/api/v1/system.go中form.PageNum
  此值还要实例化到system_service.SystemInfo中,并在其中根据此值判断是获取全部数据还是获取指定页的数据,进而从redis查询数据
*/
func GetPage(c *gin.Context) int {
	result := 0
	pageNum, _ := com.StrTo(c.Query("PageNum")).Int()
	if pageNum > 0 {
		result = (pageNum - 1)
	} else {
		return -1
	}
	return result
}
