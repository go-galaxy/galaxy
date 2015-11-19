package dao

import (
	"fmt"

	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.UserRole))
}

// Addmodel.UserRole insert a new model.UserRole into database and returns
// last inserted Id on success.
func AddUserRole(m *model.UserRole) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	id, err = o.Insert(m)
	if err == nil {
		fmt.Println("id of records insert in database:", id)
		cache.UserRoleCache.Add(m.UserId, m.RoleId)
	}
	return
}

// Getmodel.UserRoleById retrieves model.UserRole by Id. Returns error if
// Id doesn't exist
func GetUserRoleById(id int64) (v *model.UserRole, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.UserRole{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetUserRoleByUserId(id int64) (v *model.UserRole, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.UserRole{UserId: id}
	if err = o.Read(v, "user_id"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllmodel.UserRole retrieves all model.UserRole matches certain condition. Returns empty list if
// no records exist
func GetAllUserRole() (ml []*model.UserRole, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	qs := o.QueryTable("user_role")
	_, err = qs.All(&ml)
	return
}

// Updatemodel.UserRole updates model.UserRole by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserRoleById(m *model.UserRole) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.UserRole{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
			cache.UserRoleCache.Add(m.UserId, m.RoleId)
		}
	}
	return
}

// Deletemodel.UserRole deletes model.UserRole by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserRole(id int64) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.UserRole{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&model.UserRole{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
