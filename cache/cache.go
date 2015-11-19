/*
   缓存模块
   数据库修改实时修改缓存
   对外接口均为缓存数据
*/

package cache

import (
	"errors"
	"fmt"
	"github.com/go-galaxy/galaxy/model"
	"sync"
)

var (
	Users               map[int64]*model.User       //用户缓存
	Roles               map[int64]*model.Role       //角色缓存
	Permissions         map[int64]*model.Permission //权限缓存
	UserRoleCache       *UserRole                   //用户角色缓存 TODO 用户对应多个角色
	RolePermissionCache *RolePermission             //角色权限缓存
)

func Init() {
	Users = make(map[int64]*model.User)
	Roles = make(map[int64]*model.Role)
	Permissions = make(map[int64]*model.Permission)
	UserRoleCache = NewUserRole()
	RolePermissionCache = NewRolePermission()
}

type UserRole struct {
	Data map[int64]*model.Role
	Mu   *sync.Mutex
}

type RolePermission struct {
	Data map[int64][]*model.Permission
	Mu   *sync.Mutex
}

func AddUser(user *model.User) {
	Users[user.Id] = user
}

func AddRoles(m *model.Role) {
	Roles[m.Id] = m
}

func AddPermissions(user *model.Permission) {
	Permissions[user.Id] = user
}

func DelUser(id int64) {
	delete(Users, id)
}

func DelRoles(id int64) {
	delete(Roles, id)
}

func DelPermissions(id int64) {
	delete(Permissions, id)
}

func NewUserRole() *UserRole {
	return &UserRole{
		Data: make(map[int64]*model.Role),
		Mu:   &sync.Mutex{},
	}
}

func NewRolePermission() *RolePermission {
	return &RolePermission{
		Data: make(map[int64]([]*model.Permission)),
		Mu:   &sync.Mutex{},
	}
}

func (this *UserRole) Add(userId, roleId int64) (err error) {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	var role *model.Role
	if _, ok := Users[userId]; !ok {
		err = errors.New(fmt.Sprintf("没有此用户id %d", userId))
		return
	}
	if v, ok := Roles[roleId]; !ok {
		err = errors.New(fmt.Sprintf("没有此权限id %d", roleId))
		return
	} else {
		role = v
	}
	this.Data[userId] = role
	return
}

func (this *RolePermission) Add(roleId, permissionId int64) (err error) {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	var permission *model.Permission
	if v, ok := Permissions[permissionId]; !ok {
		err = errors.New(fmt.Sprintf("没有此权限id %d", permissionId))
		return
	} else {
		permission = v
	}
	if _, ok := Roles[roleId]; !ok {
		err = errors.New(fmt.Sprintf("没有此角色id %d", roleId))
		return
	}

	if v, ok := this.Data[roleId]; !ok {
		tmp := make([]*model.Permission, 0)
		tmp = append(tmp, permission)
		this.Data[roleId] = tmp
	} else {
		for _, p := range v {
			if p.Path == permission.Path && permission.Path != "" && permission.Path != "/" {
				return
			}
		}
		v = append(v, permission)
		this.Data[roleId] = v
	}
	return
}

func (this *RolePermission) Del(roleId, permissionId int64) (err error) {
	this.Mu.Lock()
	defer this.Mu.Unlock()

	if _, ok := Permissions[permissionId]; !ok {
		err = errors.New(fmt.Sprintf("没有此权限id %d", permissionId))
		return
	}
	if _, ok := Roles[roleId]; !ok {
		err = errors.New(fmt.Sprintf("没有此角色id %d", roleId))
		return
	}

	if v, ok := this.Data[roleId]; !ok {
		return
	} else {
		for k, p := range v {
			if p.Id == permissionId {
				this.Data[roleId] = append(v[0:k], v[k+1:]...)
				return
			}
		}
	}
	return
}

func GetRolePermissions(roleId int64) (p []*model.Permission) {
	p, ok := RolePermissionCache.Data[roleId]
	if !ok {
		return
	}
	return
}

func PermissionsIsExist(id int64, p []*model.Permission) bool {
	for _, v := range p {
		if v.Id == id {
			return true
		}
	}
	return false
}
