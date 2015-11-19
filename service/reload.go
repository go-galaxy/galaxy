package service

import (
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/dao"
)

//从数据库读入缓存
func Reload() (err error) {
	//reload user
	userList, err := dao.GetAllUser()
	if err != nil {
		Logs.Error("dao.GetAllUser err (%v)", err)
		return
	}
	for _, v := range userList {
		cache.Users[v.Id] = v
	}
	//reload permission
	pList, err := dao.GetAllPermission()
	if err != nil {
		Logs.Error("dao.GetAllPermission err (%v)", err)
		return
	}
	for _, v := range pList {
		cache.Permissions[v.Id] = v
	}
	//reload role
	rList, err := dao.GetAllRole()
	if err != nil {
		Logs.Error("dao.GetAllRole err (%v)", err)
		return
	}
	for _, v := range rList {
		cache.Roles[v.Id] = v
	}
	//reload userRole
	urList, err := dao.GetAllUserRole()
	if err != nil {
		Logs.Error("dao.GetAllUserRole err (%v)", err)
		return
	}
	for _, v := range urList {
		cache.UserRoleCache.Add(v.UserId, v.RoleId)
	}
	//reload rolePermission
	rpList, err := dao.GetAllRolePermission()
	if err != nil {
		Logs.Error("dao.GetAllRolePermission err (%v)", err)
		return
	}
	for _, v := range rpList {
		cache.RolePermissionCache.Add(v.RoleId, v.PermissionId)
	}
	return

}
