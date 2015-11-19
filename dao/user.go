package dao

import (
	"fmt"

	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(model.User))
}

// Addmodel.User insert a new model.User into database and returns
// last inserted Id on success.
func AddUser(m *model.User) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	m.Ctime = time.Now().Unix()
	id, err = o.Insert(m)
	if err == nil {
		cache.AddUser(m)
	}
	return
}

// Getmodel.UserById retrieves model.User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *model.User, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//Only The Superadmin Use
func GetAdminInfoByName(name, password string) (v *model.User, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.User{}
	err = o.QueryTable("user").Filter("account", name).Filter("password", password).Filter("status", model.UserStatusAdmin).One(v)
	return

}

//Get model.User by name,password,return model.User
//Return error if name or password is wrong
func GetUserInfoByName(name, password string) (v *model.User, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.User{}
	err = o.QueryTable("user").Filter("account", name).Filter("password", password).Exclude("status", model.UserStatusFreeze).One(v)
	return
}

// GetAllmodel.User retrieves all model.User matches certain condition. Returns empty list if
// no records exist
func GetAllUser() (ml []*model.User, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	qs := o.QueryTable("user").OrderBy("-ctime")
	_, err = qs.All(&ml)
	return
}

// Updatemodel.User updates model.User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *model.User) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.User{Id: m.Id}
	// ascertain id exists in the database
	err = o.Read(&v)
	if err != nil {
		return
	}
	m.Password = v.Password
	m.Ctime = time.Now().Unix()
	if num, err := o.Update(m); err == nil {
		fmt.Println("Number of records updated in database:", num)
		//修改缓存
		cache.AddUser(m)
	}
	return
}

// Deletemodel.User deletes model.User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&model.User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
			cache.DelUser(id)
		}
	}
	return
}
