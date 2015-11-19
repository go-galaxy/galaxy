package service

import (
	"errors"
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
    "github.com/astaxie/beego/logs"
)

var (
	Logs      *logs.BeeLogger
	isDebug   bool
	AdminList map[string]string
)

var (
	ErrorUserNotExist           = errors.New("没有此用户id")
	ErrorUserStatus             = errors.New("用户状态异常")
	ErrosRoleStatus             = errors.New("角色状态异常")
	ErrorUserRoleNotExist       = errors.New("用户角色无此用户")
	ErrorRolePermissnonNotExist = errors.New("角色权限无此角色")
)

func Init(dbDsn string, num int, dbName string, debug bool, adminList map[string]string) (err error) {
	AdminList = adminList
    Logs = logs.NewLogger(10000)
    Logs.SetLogger("file", `{"filename":"run.log"}`)
    Logs.SetLogger("console","")
    Logs.EnableFuncCallDepth(true)
	isDebug = debug
	model.DbName = dbName
	if dbName == "" {
		model.DbName = "default"
	}
	orm.RegisterDataBase(model.DbName, "mysql", dbDsn)
	cache.Init()
	Reload()
	Logs.Debug("重载数据完成")
	return
}

func ValidatePermission(userId int64, path string) bool {
	//查询 userId 角色 合并permission 查询path
	if user, ok := cache.Users[userId]; !ok {
		Logs.Error("err(%v) %d", ErrorUserNotExist, userId)
		return false
	} else {
		if !model.CheckUserStatus(user.Status) {
			Logs.Error("err(%v) user(%d)status(d%)", ErrorUserStatus, userId, user.Status)
			return false
		}
	}

	if v, ok := cache.UserRoleCache.Data[userId]; !ok {
		Logs.Error("err(%v) %d", ErrorUserRoleNotExist, userId)
		return false
	} else {
		if v.Status != 0 {
			return false
		}
		if p, iok := cache.RolePermissionCache.Data[v.Id]; !iok {
			Logs.Error("err(%v)  %d", ErrorRolePermissnonNotExist, v.Id)
			return false
		} else {
			//TODO  查询 permission 速度
			for _, value := range p {
				if value.Path == path && value.Status == 0 {
					return true
				}
			}
		}
	}
	return false
}

//根据用户ID获得所有菜单 缓存读取
func GetPermissionByUserId(userId int64) (mp []*model.Permission, err error) {
	if user, ok := cache.Users[userId]; !ok {
		err = ErrorUserNotExist
		return
	} else {
		if !model.CheckUserStatus(user.Status) {
			err = ErrorUserStatus
			return
		}
	}
	v, ok := cache.UserRoleCache.Data[userId]
	if !ok {
		err = ErrorUserRoleNotExist
		return
	}
	if v.Status != 0 {
		err = ErrosRoleStatus
		return
	}
	mp, ok = cache.RolePermissionCache.Data[v.Id]
	if !ok {
		err = ErrorRolePermissnonNotExist
		return
	}
	return
}
