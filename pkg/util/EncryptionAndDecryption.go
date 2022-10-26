package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type newForm struct {
	Username  string
	Password  string
	IPAddress string
}

//接收任意参数,最终返回切片类型
//func Ead(form interface{}) ([]newForm, bool) {
//	key := "1235*567@9012f4D8888$88823122234"
//	val, ok := isSlice(form)
//	if !ok {
//		return nil, false
//	}
//	//通过val.Len()获取反射值的长度
//	sliceLen := val.Len()
//	//创建一个新的切片,类型为结构体切片
//	newSlice := make([]newForm, sliceLen)
//	for i := 0; i < sliceLen; i++ {
//		//获取val值的第i个元素
//		value := val.Index(i)
//		type1 := value.Type()
//		if type1.Kind() != reflect.Struct {
//			panic("It is not struct")
//		}
//		newSlice[i].Username = value.FieldByName("Username").String()
//		newSlice[i].Password = value.FieldByName("Password").String()
//		newSlice[i].IPAddress = value.FieldByName("IPAddress").String()
//	}
//	for k, v := range newSlice {
//		EncryptCode := AesEncrypt(v.Password, key)
//		newSlice[k].Password = EncryptCode
//	}
//	return newSlice, true
//}
//
////判断是否为切片,通过反射
//func isSlice(arg interface{}) (val reflect.Value, ok bool) {
//	//使用reflect.ValudOf获取反射值对象reflect.Value
//	val = reflect.ValueOf(arg)
//	//判断val的种类是不是切片
//	if val.Kind() == reflect.Slice {
//		ok = true
//	}
//	return
//}

//加密
func AesEncrypt(str string) string {
	// 转成字节数组
	strData := []byte(str)
	k := []byte("1235*567@9012f4D8888$88823122234")
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	strData = PKCS7Padding(strData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(strData))
	// 加密
	blockMode.CryptBlocks(cryted, strData)
	return base64.StdEncoding.EncodeToString(cryted)
}

//解密
func AesDecrypt(cryted string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte("1235*567@9012f4D8888$88823122234")
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	str := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(str, crytedByte)
	// 去补全码
	str = PKCS7UnPadding(str)
	return string(str)
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(strData []byte) []byte {
	length := len(strData)
	unpadding := int(strData[length-1])
	return strData[:(length - unpadding)]
}
