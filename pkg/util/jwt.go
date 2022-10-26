package util

import (
	"oamp/global"
	"oamp/pkg/setting"
	"time"

	"github.com/golang-jwt/jwt"
)

//声明变量,类型为字节切片，值为配置文件中的加密解密秘钥串
var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type jwtClaim struct {
	username string
	password string
	jwt.StandardClaims
}

//创建token
func CreateToken(username, password string) (string, error) {
	//time.Now返回当前时间,返回类型为Time结构体
	nowTime := time.Now()
	//调用Time结构体的Add()方法,在当前时间上+12小时
	expireTime := nowTime.Add(12 * time.Hour)
	claim := &jwtClaim{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gongguan",
		},
	}
	//SigningMethodHS256表示加密方法,claim表示加密参数，为结构体,返回值为tokenClaims结构体
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	//调用加密方法SignedString,参数为配置文件中定义的加密秘钥
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

//解析token
func ParseToken(token string) (*jwtClaim, error) {
	/*
	   ParseWithClaims用户解析鉴权的声明,方法内部主要是解码和校验过程,最终返回*Token
	   tokenClaim实际是一个*Token
	*/
	tokenClaim, err := jwt.ParseWithClaims(token, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		global.Log.Error(err.Error())
	}
	if tokenClaim != nil {
		claims, ok := tokenClaim.Claims.(*jwtClaim)
		if ok && tokenClaim.Valid {
			return claims, nil
		}
	}
	return nil, err
}
