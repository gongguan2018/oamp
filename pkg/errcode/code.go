package errcode

//声明错误代码常量
const (
	SUCCESS                          = 200
	ERROR_IP_EXISTS                  = 10001
	ERROR_IP_NOT_EXISTS              = 10002
	ERR_USER_NOT_EXISTS              = 10003
	ERROR                            = 10004
	INVALID_PARAMS                   = 10005
	ERROR_PASS                       = 10006
	LOGIN_SUCCESS                    = 10010
	ERROR_AUTH_CHECK_TOKEN_FAIL      = 10011
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT   = 10012
	ERROR_PARAMS_BIND_FAIL           = 10013
	IP_EXISTS                        = 10014
	ID_NOT_EXISTS                    = 10015
	ERROR_GET_DATA_FAIL              = 10016
	ERROR_GET_USERPASS_FAIL          = 10017
	INVALID_PARAMS_USER_PASS_NOT_NIL = 10018
	USER_EXISTS                      = 10019
	DEFAULT_USER_CANNOT_BE_DELETE    = 10020
)
