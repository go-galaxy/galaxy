package model

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	DbName string
)

const (
	UserStatusAdmin  = 2 //管理员状态
	UserStatusFreeze = 1 //冻结状态
	UserStatusNormal = 0 //正常状态
)

func Md5(str string) string {
	newStr := md5.New()
	newStr.Write([]byte(str))
	return hex.EncodeToString(newStr.Sum(nil))
}

func CheckUserStatus(status int) bool {
	if status == UserStatusNormal || status == UserStatusAdmin {
		return true
	}
	return false
}
