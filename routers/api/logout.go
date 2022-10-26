package api

import (
	"net/http"
	"oamp/pkg/app"
	"oamp/pkg/errcode"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	var response = app.Gin{C: c}
	data := make(map[string]interface{})
	data["logout"] = "SUCCESS"
	response.Response(http.StatusOK, errcode.SUCCESS, data)
}
