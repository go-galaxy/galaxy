package routers

import (
	"github.com/go-galaxy/galaxy/controller"
	"github.com/go-martini/martini"
)

func Router(m *martini.ClassicMartini) {
	m.Get("/", controller.Index)
	m.Get("/role", controller.Role)
	m.Get("/permission", controller.Permission)
	m.Get("/permission/add", controller.PermissonAdd)
	m.Post("/permission/add", controller.PermissonAdd)
	m.Get("/role/add", controller.RoleAdd)
	m.Post("/role/add", controller.RoleAdd)
	m.Post("/role_permission/add", controller.RolePermissionAdd)
	m.Get("/user/edit", controller.UserEdit)
	m.Post("/user/edit", controller.UserEdit)
	m.Post("/user/add", controller.UserAdd)
	m.Get("/user/add", controller.UserAdd)
	m.Post("/user_role/edit", controller.UserRoleEdit)
	m.Get("/test", controller.UserPermissionTest)
	m.Get("/login", controller.Login)
	m.Post("/login", controller.Login)
	m.Get("/logout", controller.Logout)
	m.Get("/create/admin", controller.CreateAdmin)
	m.Post("/create/admin", controller.CreateAdmin)
}
