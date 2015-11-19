package dao

import (
	"fmt"

	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.Role))
}

// Addmodel.Role insert a new model.Role into database and returns
// last inserted Id on success.
func AddRole(m *model.Role) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	id, err = o.Insert(m)
	if err == nil {
		m.Id = id
		cache.Roles[id] = m
	}
	return
}

// Getmodel.RoleById retrieves model.Role by Id. Returns error if
// Id doesn't exist
func GetRoleById(id int64) (v *model.Role, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.Role{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllmodel.Role retrieves all model.Role matches certain condition. Returns empty list if
// no records exist
func GetAllRole() (ml []*model.Role, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	qs := o.QueryTable("role")
	_, err = qs.All(&ml)
	return
}

// Updatemodel.Role updates model.Role by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoleById(m *model.Role) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
			cache.Roles[m.Id] = m
		}
	}
	return
}

// Deletemodel.Role deletes model.Role by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRole(id int64) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&model.Role{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
			delete(cache.Roles, id)
		}
	}
	return
}
