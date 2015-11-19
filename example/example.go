package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-galaxy/galaxy"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strings"
)

var (
	ga *galaxy.Galaxy
)

func init() {
	ga = &galaxy.Galaxy{}
}

func main() {
	ga.AdminPort = 8000
	ga.MysqlDsn = "root:scjh1234@tcp(192.168.1.15:3306)/galaxy2"
	ga.TemplatePath = "/Users/sunlei/go/src/galaxy/templates"
	ga.DbName = "default"
	go ga.Run()
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(render.Renderer(render.Options{
		Directory:  "/Users/sunlei/go/src/github.com/go-galaxy/galaxy/example/template", // Specify what path to load the templates from.
		Extensions: []string{".tmpl", ".html"},                     // Specify extensions to load for templates.
		Layout:     "",
	}))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(func(res http.ResponseWriter, req *http.Request, session sessions.Session) {
		if strings.Contains(req.URL.String(), "login") {
			return
		}
		v := session.Get("login")
		if v == nil {
			http.Redirect(res, req, "login", http.StatusFound)
		}
	})
	m.Use(sessions.Sessions("my_session", store))
	m.Get("/login", Login)
	m.Get("/", Index)
	m.Post("/login", Login)
	m.Get("/logout", LogOut)
	m.Get("/test", Test)
	m.RunOnAddr(":3000")
}

func Md5(str string) string {
	newStr := md5.New()
	newStr.Write([]byte(str))
	return hex.EncodeToString(newStr.Sum(nil))
}

func Index(session sessions.Session, r render.Render) {
	id, ok := session.Get("user_id").(int64)
	if !ok {
		r.Data(200, []byte("error"))
		return
	}

	list, err := ga.GetPermissionByUserId(id)
	if err != nil {
		r.Data(200, []byte("get error"))
		return
	}

	r.HTML(200, "list", map[string]interface{}{"list": list})
	return
}

func Test(session sessions.Session, req *http.Request) string {
	req.ParseForm()
	values := req.Form
	path := values.Get("path")

	id, ok := session.Get("user_id").(int64)
	if !ok {
		return "error"
	}
	ok = ga.ValidatePermission(id, path)
	return fmt.Sprintf("%v", ok)
}

func Login(req *http.Request, session sessions.Session) string {
	h := `
    <!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, height=device-height, initial-scale=0.0, maximum-scale=0.0, minimum-scale=0.4, user-scalable=0"/>
    <meta name="format-detection" content="telephone=no"/>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <title>galaxy权限系统</title>
    <!-- 新 Bootstrap 核心 CSS 文件 -->
    <link rel="stylesheet" href="http://cdn.bootcss.com/bootstrap/3.3.4/css/bootstrap.min.css">

    <!-- 可选的Bootstrap主题文件（一般不用引入） -->
    <link rel="stylesheet" href="http://cdn.bootcss.com/bootstrap/3.3.4/css/bootstrap-theme.min.css">


</head>
<body style="margin: 0 20px">
<nav class="navbar navbar-default">
    <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="/">主页</a>

            <ul class="nav navbar-nav">
                <li class="active"><a href="/logout">登出</a></li>
            </ul>
        </div>
    </div><!-- /.container-fluid -->
</nav>
<!-- jQuery文件。务必在bootstrap.min.js 之前引入 -->
<script src="http://cdn.bootcss.com/jquery/1.11.2/jquery.min.js"></script>

<!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
<script src="http://cdn.bootcss.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
 <form id="loginForm" method="" action="">
        <div class="form-group">
            <label for="account">账号</label>
            <input type="text" class="form-control" id="account" name="account" placeholder="账号">
        </div>
        <div class="form-group">
            <label for="pass">密码</label>
            <input type="password" class="form-control" id="pass" name="pass" placeholder="密码">
        </div>

        <button type="button" class="btn btn-default" onclick="login()" >登陆</button>
    </form>
    <script>

        function login(){
            var account=$("#account").val();
            if (account==""){
                alert("账号不能为空！");
                return false;
            }

            var password=$("#pass").val();
            if (password==""){
                alert("密码不能为空！");
                return false;
            }
             $.post("/login",
                    {
                        account: account,
                        pass: password,
                    },
                    function(data, status) {
                        if (data==0){
                            window.location.href="/";
                        }else{
                            alert(data);
                        }
                    }
            );
        }
    </script>
<script>

</script>
</body>
</html>

    `
	if req.Method == "GET" {
		return h
	}
	req.ParseForm()
	values := req.Form
	account := values.Get("account")
	pass := Md5(values.Get("pass"))

	//returnData["msg"] = "登录失败"
	var errMsg string
	if account == "" || pass == "" {
		errMsg = "pass 或者  account 为空"
		return errMsg
	}
	userInfo, err := ga.GetUserInfo(account, pass)
	if err != nil {
		fmt.Println(err)
		errMsg = "登陆失败,账号或密码错误,或者账号已经禁用！"
		return errMsg
	}

	session.Set("login", "true")
	session.Set("user_id", userInfo.Id)
	session.Set("user_name", userInfo.Name)
	return "0"
}

func LogOut(session sessions.Session) string {
	session.Clear()
	return "退出成功"
}
