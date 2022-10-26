package api

import (
	"net/http"
	"oamp/pkg/app"
	"oamp/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type Info struct {
	Username string `json:"username"`
	Role     string `json:"superadmin"`
	Avatar   string `json:"avatar"`
}

func UserInfo(c *gin.Context) {
	var (
		admin      = "admin"
		superadmin = "superadmin"
		avatar     = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	)
	appG := app.Gin{C: c}
	info := Info{
		Username: admin,
		Role:     superadmin,
		Avatar:   avatar,
	}
	data := make(map[string]interface{})
	data["username"] = info.Username
	data["Role"] = info.Role
	data["avatar"] = info.Avatar
	appG.Response(http.StatusOK, errcode.SUCCESS, data)
}
