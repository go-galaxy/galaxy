package service

import (
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/dao"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
	"html/template"
	"time"
	//	"fmt"
)

func AddUser(id int64, name, info, pass string) (err error) {
	m := &model.User{}
	m.Id = id
	m.Name = name
	m.Info = info
	m.Ctime = time.Now().Unix()
	_, err = dao.AddUser(m)
	return
}

func DelUser(id int64) (err error) {
	err = dao.DeleteUser(id)
	return
}

func EditUser(id int64, name, info string, status int, pass string) (err error) {
	m, err := dao.GetUserById(id)
	if err != nil {
		Logs.Error("dao.GetUserById err(%v)", err)
		return
	}
	m.Name = name
	m.Info = info
	m.Status = status
	err = dao.UpdateUserById(m)
	return
}

func EditUserRole(userId, roleId int64) (err error) {
	_, err = dao.GetUserById(userId)
	if err != nil {
		Logs.Error("dao getuserbudi err(%v)", err)
		return
	}
	_, err = dao.GetRoleById(roleId)
	if err != nil {
		Logs.Error("dao GetRoleById err(%v)", err)
		return
	}
	ur, err := dao.GetUserRoleByUserId(userId)
	if err != nil && err != orm.ErrNoRows {
		Logs.Error("dao GetUserRoleByUserId err(%v)", err)
		return
	}
	if err == orm.ErrNoRows {
		ur = &model.UserRole{}
		ur.UserId = userId
		ur.RoleId = roleId
		_, err = dao.AddUserRole(ur)
		if err != nil {
			Logs.Error("dao AddUserRole err(%v)", err)
			return
		}
		return
	}
	ur.RoleId = roleId
	err = dao.UpdateUserRoleById(ur)
	if err != nil {
		Logs.Error("dao UpdateUserRoleById err(%v)", err)
		return
	}
	return
}

func GetUserRole(in int64) *model.Role {
	out := &model.Role{}
	v, ok := cache.UserRoleCache.Data[in]
	if !ok {
		//Logs.Error("dao.GetUserRoleByUserId err %v", err)
		return out
	}
	return v
}

//根据用户名和密码验证是否可以登录 从数据库读取
func GetUserInfoByName(name, password string) (out *model.User, err error) {
	out = &model.User{}
	out, err = dao.GetUserInfoByName(name, password)
	if err != nil {
		return
	}
	return
}

//根据用户ID获得所有菜单
//func GetPermissionByRoleId(role_id int64) ([]*model.Permission, error) {
//	return dao.GetPermissionByRoleId(role_id)
//}

//解析时间格式
func FormatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	//return time.Now().Format("2006-01-02 15:04:05")
}

var FuncMap template.FuncMap

func init() {
	FuncMap = make(template.FuncMap)
	FuncMap["getUserRole"] = GetUserRole
	FuncMap["formatTime"] = FormatTime
}
