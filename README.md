galaxy
===============

权限系统
例子见  example

    package main

    import (
        "galaxy"
        "time"
    )
    func main() {
        a := &galaxy.Galaxy{}
        a.AdminPort = 8000
        a.MysqlDsn = "user:pass@tcp(127.0.0.1:3306)/galaxy"
        a.TemplatePath = "/path/galaxy/templates"
        //如果已经存在beego orm 请修改default为其他
        a.DbName = "default"
        go a.Run()
        time.Sleep(10000000 * time.Second)
    }

调用方式 以beego为例 在router.go中

    var FilterUser = func(ctx *context.Context) {

    	if ctx.Request.RequestURI == "/Login" {
    		return
    	}

    	if ctx.Request.RequestURI == "/LoginOut" {
    		return
    	}

    	if strings.HasPrefix(ctx.Request.RequestURI, "/static") {
    		return
    	}

    	if admin, ok := ctx.Input.CruSession.Get("adminInfo").(*models.Admin); !ok {
    		ctx.WriteString(`{"Status":"2","ErrInfo":"您还没有登录！请您先登录！"}`)
    	} else {
    		if ok := galaxy.ValidatePermission(admin.UserId, ctx.Request.RequestURI); !ok {
    		    ctx.WriteString(`{"Status":"0","ErrInfo":"您没有权限访问！"}`)
    		}
    	}

    }
    func init() {
        beego.Router("/", &controllers.Index{}) //根目录
        beego.InsertFilter("*", beego.BeforeRouter, FilterUser)
    }

自己完成注册用户可以调用
    galaxy.AddUser galaxy.EditUser galaxy.DelUser

登录验证
    galaxy.GetUserInfo(name, password)

导入数据库
    首先需要创建数据库
    或者导入到已有数据库 注意会覆盖原有数据库 有user表的注意了
    mysqldump -u用户名 -p密码 -h主机 数据库 < doc/galaxy.sql

要完善的功能：
    修改管理员密码