/*
	权限控制
    存储落地默认mysql
    web 控制页面
    权限
    群组
    用户
*/
package galaxy

import (
	"errors"
	"fmt"
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/go-galaxy/galaxy/routers"
	"github.com/go-galaxy/galaxy/service"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Galaxy struct {
	AdminPort    int    //管理界面端口
	MysqlDsn     string //user:pass@tcp(127.0.0.1:3306)/galaxy?timeout=3s&parseTime=true&loc=Local&charset=utf8mb4,utf8
	MysqlNum     int
	TemplatePath string
	Debug        bool
	DbName       string
	AdminList    map[string]string
}

func (this *Galaxy) Run() (err error) {
	if this.TemplatePath == "" || this.MysqlDsn == "" || this.AdminPort < 3000 {
		return fmt.Errorf("参数错误 请检查参数")
	}
	service.Init(this.MysqlDsn, this.MysqlNum, this.DbName, this.Debug, this.AdminList)
	//运行管理界面
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  this.TemplatePath,          // Specify what path to load the templates from.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Layout:     "layout",
		Funcs:      []template.FuncMap{service.FuncMap},
	}))
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(func(res http.ResponseWriter, req *http.Request, session sessions.Session) {
		if strings.Contains(req.URL.String(), "login") {
			return
		}
		if strings.HasPrefix(req.URL.String(), "/create/admin") {
			return
		}
		if len(cache.Users) == 0 {
			http.Redirect(res, req, "/create/admin", http.StatusFound)
			return
		}
		v := session.Get("login")
		if v == nil {
			http.Redirect(res, req, "login", http.StatusFound)
		}
	})
	routers.Router(m)
	service.Logs.Debug("listen on :%d", this.AdminPort)
	m.RunOnAddr(fmt.Sprintf(":%d", this.AdminPort))
	return
}

//添加用户 pass为加密之后的字符串
func (this *Galaxy) AddUser(id int64, name, info, pass string) (err error) {
	if name == "" || info == "" {
		return errors.New("args err")
	}
	return service.AddUser(id, name, info, pass)
}

//删除用户
func (this *Galaxy) DelUser(id int64) (err error) {
	return service.DelUser(id)
}

//编辑用户 pass为加密之后的字符串
func (this *Galaxy) EditUser(id int64, name, info string, status int, pass string) (err error) {
	if name == "" || info == "" || (status != 0 || status != 1) {
		return errors.New("args err")
	}
	return service.EditUser(id, name, info, status, pass)
}

//验证权限
func (this *Galaxy) ValidatePermission(userId int64, path string) bool {
	return service.ValidatePermission(userId, path)
}

//获取用户详情 pass为加密之后的字符串
func (this *Galaxy) GetUserInfo(name, password string) (*model.User, error) {
	if name == "" || password == "" {
		return nil, errors.New("args err")
	}
	return service.GetUserInfoByName(name, password)
}

//获取角色ID
func (this *Galaxy) GetUserRole(userId int64) *model.Role {
	return service.GetUserRole(userId)
}

//根据角色ID获取权限列表
func (this *Galaxy) GetPermissionByUserId(userId int64) (mp []*model.Permission, err error) {
	mp, err = service.GetPermissionByUserId(userId)
	if err != nil {
		service.Logs.Error("GetPermissionByUserId userId(%d) err(%v)", userId, err)
	}
	return
}

//写入操作日志
func (this *Galaxy) AddLog(user_id int64, role_id int, logContent, tag, ipv4 string) error {
	mylog := &model.ActionLog{}
	mylog.UserId = user_id
	mylog.RoleId = role_id
	mylog.LogContent = logContent
	mylog.Tag = tag
	mylog.Ipv4 = ipv4
	mylog.CreateTime = time.Now().Unix()
	fmt.Println(mylog)
	return service.AddLog(mylog)
}
