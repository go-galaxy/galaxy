package dao

import (
	"fmt"
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.RolePermission))
}

// Addmodel.RolePermission insert a new model.RolePermission into database and returns
// last inserted Id on success.
func AddRolePermission(m *model.RolePermission) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	id, err = o.Insert(m)
	if err == nil {
		err = cache.RolePermissionCache.Add(m.RoleId, m.PermissionId)
		fmt.Println(err)
	}
	return
}

// Getmodel.RolePermissionById retrieves model.RolePermission by Id. Returns error if
// Id doesn't exist
func GetRolePermissionById(id int64) (v *model.RolePermission, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.RolePermission{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllmodel.RolePermission retrieves all model.RolePermission matches certain condition. Returns empty list if
// no records exist
func GetAllRolePermission() (ml []*model.RolePermission, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	qs := o.QueryTable("role_permission")
	_, err = qs.All(&ml)
	return
}

// Updatemodel.RolePermission updates model.RolePermission by Id and returns error if
// the record to be updated doesn't exist
func UpdateRolePermissionById(m *model.RolePermission) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.RolePermission{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
			cache.RolePermissionCache.Add(m.RoleId, m.PermissionId)
		}
	}
	return
}

// Deletemodel.RolePermission deletes model.RolePermission by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRolePermission(id int64) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.RolePermission{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&model.RolePermission{Id: id}); err == nil {
			cache.RolePermissionCache.Del(v.RoleId, v.PermissionId)
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//根据
func DeleteRoleIdAndPermissionId(role_id int64, permission_id int64) error {
	o := orm.NewOrm()
	o.Using(model.DbName)
	//	_, err := o.Raw("DELETE FROM `role_permission` WHERE role_id= ? and permission_id = ?", role_id, permission_id).Exec()
	//	if err != nil {
	//		return err
	//	}
	if num, err := o.Delete(&model.RolePermission{RoleId: role_id, PermissionId: permission_id}); err == nil {
		fmt.Println(num)
	}

	cache.RolePermissionCache.Del(role_id, permission_id)
	//fmt.Println("Number of records deleted in database:")
	return nil
}
