package model

type Permission struct {
	Id       int64  `orm:"column(id);auto"`
	Name     string `orm:"column(name);size(45)"`
	Info     string `orm:"column(info);size(100)"`
	Path     string `orm:"column(path);size(100)"`
	Status   int    `orm:"column(status)"`
	Sort     int    `orm:"column(sort);"`
	ParentId int    `orm:"column(parent_id);"`
	Types    int    `orm:"column(types);"`
}

type Role struct {
	Id     int64  `orm:"column(id);auto"`
	Name   string `orm:"column(name);size(100)"`
	Info   string `orm:"column(info);size(100)"`
	Status int    `orm:"column(status)"`
}

type RolePermission struct {
	Id           int64 `orm:"column(id);pk"`
	RoleId       int64 `orm:"column(role_id)"`
	PermissionId int64 `orm:"column(permission_id)"`
}

type User struct {
	Id       int64  `orm:"column(id);pk"`
	Account  string `orm:"column(account);size(100)"`
	Password string `orm:"column(password);size(50)"`
	Name     string `orm:"column(name);size(45)"`
	Info     string `orm:"column(info);size(45)"`
	Status   int    `orm:"column(status)"`
	Ctime    int64  `orm:"column(ctime)"`
	//RoleId   int    `orm:"column(role_id)"`
}

type UserRole struct {
	Id     int64 `orm:"column(id);pk"`
	UserId int64 `orm:"column(user_id)"`
	RoleId int64 `orm:"column(role_id)"`
}

type ActionLog struct {
	Id         int    `orm:"column(id);pk"`
	UserId     int64  `orm:"column(userId)"`
	RoleId     int    `orm:"column(roleId)"`
	CreateTime int64  `orm:"column(createTime)"`
	LogContent string `orm:"column(logContent)"`
	Tag        string `orm:"column(tag)"`
	Ipv4       string `orm:"column(ipv4)"`
}

type ShowPermission struct {
	ParentInfo Permission
	SubInfo    []Permission
}
