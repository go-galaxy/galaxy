/*
   管理后台控制器
*/

package controller

import (
	"fmt"
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/dao"
	"github.com/go-galaxy/galaxy/model"
	"github.com/go-galaxy/galaxy/service"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strconv"
	"strings"
)

func Init() (err error) {
	return
}

func Index(r render.Render) {
	users, err := dao.GetAllUser()
	if err != nil {
		service.Logs.Error("dao.GetAllUser() err(%v)", err)
		return
	}
	data := make(map[string]interface{})
	rs, err := dao.GetAllRole()
	if err != nil {
		service.Logs.Error("dao.GetAllUser() err(%v)", err)
		return
	}
	data["users"] = users
	data["r"] = rs
	r.HTML(200, "admin", data)
	return
}

func Role(r render.Render) {
	l, err := dao.GetAllRole()
	if err != nil {
		service.Logs.Error("dao.GetAllRole() err(%v)", err)
	}
	data := make(map[string]interface{})
	data["l"] = l
	r.HTML(200, "role", data)
	return
}

func Permission(r render.Render) {
	l, err := dao.GetAllPermission()
	if err != nil {
		service.Logs.Error("dao.GetAllPermission() err(%v)", err)
	}
	data := make(map[string]interface{})
	data["l"] = l
	r.HTML(200, "permission", data)
	return
}

func RoleAdd(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	id := values.Get("id")
	var idInt64 int64
	var err error
	m := &model.Role{}
	data := make(map[string]interface{})

	l, err := dao.GetAllPermission()
	if err != nil {
		service.Logs.Error("dao.GetAllPermission() err(%v)", err)
	}
	data["l"] = l

	list := dao.GetFormatPermission()
	data["list"] = list

	if id != "" {
		idInt64, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			service.Logs.Error("strconv.ParseInt err(%v)", err)
			return
		}
	}
	if req.Method == "GET" {
		if id != "" && idInt64 != 0 {
			m, err = dao.GetRoleById(idInt64)
			if err != nil {
				service.Logs.Error("dao.GetRoleById err(%v)", err)
				return
			}
		}
		p := cache.GetRolePermissions(idInt64)
		data["p"] = p
		data["m"] = m
		r.HTML(200, "role_add", data)
		return
	}
	m.Info = values.Get("info")
	m.Name = values.Get("name")
	if id == "" || idInt64 == 0 {
		m.Status = 0
		mId, err := dao.AddRole(m)
		if err != nil {
			service.Logs.Error("dao.AddRole err(%v)", err)
			return
		}
		m, _ = dao.GetRoleById(mId)
	} else {
		m.Id = idInt64
		status, _ := strconv.Atoi(values.Get("status"))
		m.Status = status
		err = dao.UpdateRoleById(m)
		if err != nil {
			service.Logs.Error("dao.UpdateRoleById err(%v)", err)
			return
		}
	}
	p := cache.GetRolePermissions(idInt64)

	data["p"] = p
	data["m"] = m

	r.HTML(200, "role_add", data)
	return
}

func PermissonAdd(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	id := values.Get("id")
	var idInt64 int64
	var err error
	m := &model.Permission{}
	data := make(map[string]interface{})
	if id != "" {
		idInt64, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			service.Logs.Error("strconv.ParseInt err(%v)", err)
			return
		}
	}
	if req.Method == "GET" {
		if id != "" && idInt64 != 0 {
			m, err = dao.GetPermissionById(idInt64)
			if err != nil {
				service.Logs.Error("dao.GetPermissionById err(%v)", err)
				return
			}
		}
		data["m"] = m
		r.HTML(200, "permission_add", data)
		return
	}
	m.Info = values.Get("info")
	m.Name = values.Get("name")
	m.Path = values.Get("path")
	if id == "" || idInt64 == 0 {
		m.Status = 0
		mId, err := dao.AddPermission(m)
		if err != nil {
			service.Logs.Error("dao.AddPermission err(%v)", err)
			return
		}
		m, _ = dao.GetPermissionById(mId)
	} else {
		m.Id = idInt64
		status, _ := strconv.Atoi(values.Get("status"))
		m.Status = status
		err = dao.UpdatePermissionById(m)
		if err != nil {
			service.Logs.Error("dao.UpdatePermissionById err(%v)", err)
			return
		}
	}
	data["m"] = m
	r.HTML(200, "permission_add", data)
	return
}

//管理后台修改用户
func UserEdit(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	id := values.Get("id")
	var idInt64 int64
	var err error
	m := &model.User{}
	data := make(map[string]interface{})
	l, err := dao.GetAllRole()
	if err != nil {
		service.Logs.Error("dao.GetAllRole() err(%v)", err)
		return
	}
	data["l"] = l

	if id == "" {
		service.Logs.Error("id==null)")
		return
	}

	idInt64, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		service.Logs.Error("strconv.ParseInt err(%v)", err)
		return
	}

	rs, err := dao.GetAllRole()
	if err != nil {
		service.Logs.Error("dao.GetAllUser() err(%v)", err)
		return
	}
	data["roleList"] = rs

	if req.Method == "GET" {
		m, err = dao.GetUserById(idInt64)
		if err != nil {
			service.Logs.Error("dao.GetRoleById err(%v)", err)
			return
		}
		data["m"] = m
		r.HTML(200, "user_edit", data)
		return
	}
	m.Account = values.Get("account")
	m.Info = values.Get("info")
	m.Name = values.Get("name")
	m.Id = idInt64
	status, _ := strconv.Atoi(values.Get("status"))
	m.Status = status
	roleId := values.Get("role_id")
	roleIdInt, _ := strconv.Atoi(roleId)

	//更新用户角色
	err = service.EditUserRole(m.Id, int64(roleIdInt))
	if err != nil {
		service.Logs.Error("service.EditUserRole err(%v)", err)
		return
	}

	//更新用户
	err = dao.UpdateUserById(m)
	if err != nil {
		service.Logs.Error("dao.UpdateUserById err(%v)", err)
		return
	}
	data["m"] = m
	r.HTML(200, "user_edit", data)
	return
}

//后台添加用户
func UserAdd(r render.Render, req *http.Request) {
	var err error
	data := make(map[string]interface{})
	if req.Method == "GET" {
		rs, err := dao.GetAllRole()
		if err != nil {
			service.Logs.Error("dao.GetAllUser() err(%v)", err)
			return
		}
		data["roleList"] = rs
		r.HTML(200, "user_add", data)
		return
	}
	req.ParseForm()
	values := req.Form
	m := &model.User{}
	m.Account = values.Get("account")
	m.Password = values.Get("password")
	m.Info = values.Get("info")
	m.Name = values.Get("name")
	if !checkNull([]string{m.Account, m.Password, m.Info, m.Name}...) {
		service.Logs.Error("args err")
		return
	}
	status, _ := strconv.Atoi(values.Get("status"))
	m.Status = status
	roleId, _ := strconv.Atoi(values.Get("role_id"))
	//添加用户

	userId, err := dao.AddUser(m)
	if err != nil {
		service.Logs.Error("dao.InsertUser err(%v)", err)
		return
	}
	//添加用户角色
	ur := &model.UserRole{}
	ur.RoleId = int64(roleId)
	ur.UserId = userId
	_, err = dao.AddUserRole(ur)
	if err != nil {
		service.Logs.Error("dao.AddUserRole err(%v)", err)
		return
	}

	r.Redirect("/", 302)
	return
}

//管理后台 加减权限
func RolePermissionAdd(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	var err error
	//角色id
	roleId := values.Get("role_id")
	roleIdInt, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		service.Logs.Error("strconv.ParseInt err(%v)", err)
		return
	}
	//权限列表
	permissionsAddStr := values.Get("permission_add_list")
	service.Logs.Debug("add list %s", permissionsAddStr)

	existIds := cache.GetRolePermissions(roleIdInt)

	addList := make(map[int64]int64)
	for _, v := range strings.Split(permissionsAddStr, ",") {
		if v == "" {
			continue
		}
		idInt, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			addList[idInt] = idInt
		}
	}

	if len(existIds) > 0 {
		for _, n := range existIds {
			if _, ok := addList[n.Id]; ok {
				delete(addList, n.Id)

			}
		}
	}

	for _, v := range addList {
		m := &model.RolePermission{}
		m.RoleId = roleIdInt
		m.PermissionId = v
		_, err = dao.AddRolePermission(m)
		if err != nil {
			service.Logs.Error("dao.AddRolePermission err(%v)", err)
			continue
		}
	}
	permissionsDelStr := values.Get("permission_del_list")
	service.Logs.Debug("del list %s", permissionsDelStr)

	delList := make(map[int64]int64)
	for _, v := range strings.Split(permissionsDelStr, ",") {
		if v == "" {
			continue
		}
		for _, e := range existIds {
			if v == strconv.FormatInt(e.Id, 10) {
				if _, ok := delList[e.Id]; !ok {
					delList[e.Id] = e.Id
				} else {
					fmt.Println("Key IsFound")
				}
			}
		}
	}

	for _, v := range delList {
		err = dao.DeleteRoleIdAndPermissionId(roleIdInt, v)
		if err != nil {
			service.Logs.Error("dao.DelRolePermission err(%v)", err)
			continue
		}
	}
	service.Logs.Debug("redirect /role/add?id=%d", roleIdInt)
	r.Redirect(fmt.Sprintf("/role/add?id=%d", roleIdInt), 302)
	return
}

//用户 角色分配修改
func UserRoleEdit(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	var err error
	data := make(map[string]interface{})
	data["ret"] = 1
	userIdStr := values.Get("user_id")
	roleIdStr := values.Get("role_id")
	if userIdStr == "" || roleIdStr == "" {
		service.Logs.Error("args err")
		r.JSON(200, data)
		return
	}
	userIdInt, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		service.Logs.Error("strconv.ParseInt err(%v)", err)
		r.JSON(200, data)
		return
	}
	roleIdInt, err := strconv.ParseInt(roleIdStr, 10, 64)
	if err != nil {
		service.Logs.Error("strconv.ParseInt err(%v)", err)
		r.JSON(200, data)
		return
	}
	err = service.EditUserRole(userIdInt, roleIdInt)
	if err != nil {
		data["ret"] = 1
		r.JSON(200, data)
		return
	}
	data["ret"] = 0
	r.JSON(200, data)
	return
}

func UserPermissionTest(r render.Render, req *http.Request) {
	req.ParseForm()
	values := req.Form
	var err error
	data := make(map[string]interface{})
	userIdStr := values.Get("user_id")
	path := values.Get("path")
	if userIdStr == "" || path == "" {
		service.Logs.Error("args err")
		return
	}
	userIdInt, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		service.Logs.Error("strconv.ParseInt err(%v)", err)
		return
	}
	ok := service.ValidatePermission(userIdInt, path)
	data["ret"] = ok
	r.JSON(200, data)
	return
}

func Login(r render.Render, req *http.Request, session sessions.Session) {
	if req.Method == "GET" {
		r.HTML(200, "login", "")
		return
	}
	req.ParseForm()
	values := req.Form
	account := values.Get("account")
	fmt.Println(values.Get("pass"))
	pass := model.Md5(values.Get("pass"))
	fmt.Println(pass)
	fmt.Println(account)
	returnData := make(map[string]interface{})
	returnData["status"] = 0
	//returnData["msg"] = "登录失败"
	var errMsg string
	if account == "" || pass == "" {
		errMsg = "pass 或者  account 为空"
		service.Logs.Error(errMsg)
		returnData["msg"] = errMsg
		returnData["status"] = 1
		r.HTML(200, "login", returnData)
		return
	}
	userInfo, err := dao.GetAdminInfoByName(account, pass)
	if err != nil {
		errMsg = "登陆失败,账号或密码错误,或者账号已经禁用！"
		service.Logs.Error(fmt.Sprintf("%s", err))
		returnData["msg"] = errMsg
		returnData["status"] = 1
		r.JSON(200, returnData)
		return
	}

	session.Set("login", "true")
	session.Set("UserName", userInfo.Name)
	//	r.JSON(200, map[string]interface{}{"error": "10001", "msg": "Data Error"})
	//	return

	returnData["msg"] = "登陆成功！"
	r.JSON(200, returnData)
	return
}

func Logout(res http.ResponseWriter, req *http.Request, session sessions.Session) {
	session.Delete("login")
	session.Delete("UserName")
	http.Redirect(res, req, "/login", http.StatusFound)
	return
}

func checkNull(args ...string) bool {
	for _, v := range args {
		if v == "" {
			return false
		}
	}
	return true
}

func CreateAdmin(r render.Render, req *http.Request) {
	if len(cache.Users) > 0 {
		r.Error(404)
		return
	}
	if req.Method == "GET" {
		r.HTML(200, "user_add", map[string]interface{}{})
		return
	}

	req.ParseForm()
	values := req.Form
	m := &model.User{}
	m.Account = values.Get("account")
	m.Password = model.Md5(values.Get("password"))
	m.Info = values.Get("info")
	m.Name = values.Get("name")
	if !checkNull([]string{m.Account, m.Password, m.Info, m.Name}...) {
		service.Logs.Error("args err")
		return
	}
	m.Status = model.UserStatusAdmin
	_, err := dao.AddUser(m)
	if err != nil {
		service.Logs.Error("dao.InsertUser err(%v)", err)
		return
	}
	r.Redirect("/", 302)
	return
}
