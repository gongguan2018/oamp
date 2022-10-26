package errcode

//声明映射，键为错误代码数字，值为对应的字符串
//此映射的键即为code.go中声明的常量的键,同一包中的常量可以直接调用
var Msg = map[int]string{
	SUCCESS:                          "OK",
	ERROR_IP_EXISTS:                  "IP存在",
	ERROR_IP_NOT_EXISTS:              "IP不存在",
	ERR_USER_NOT_EXISTS:              "用户不存在",
	ERROR:                            "fail",
	INVALID_PARAMS:                   "请求参数错误",
	ERROR_PASS:                       "密码错误",
	LOGIN_SUCCESS:                    "登录页面成功",
	ERROR_AUTH_CHECK_TOKEN_FAIL:      "token解析失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:   "token超时",
	ERROR_PARAMS_BIND_FAIL:           "参数绑定失败",
	IP_EXISTS:                        "IP已经存在",
	ID_NOT_EXISTS:                    "ID不存在",
	ERROR_GET_DATA_FAIL:              "获取数据失败",
	ERROR_GET_USERPASS_FAIL:          "获取用户名密码失败",
	INVALID_PARAMS_USER_PASS_NOT_NIL: "请求参数错误,用户名和密码不能为空",
	USER_EXISTS:                      "用户已存在",
	DEFAULT_USER_CANNOT_BE_DELETE:    "默认用户不可删除",
}

//此函数根据code.go中的常量作为Msg中的key,并获取对应的值返回
func GetMsg(code int) string {
	msg, ok := Msg[code] //ok表示是否存在key
	if ok {
		return msg
	}
	return Msg[ERROR]
}
