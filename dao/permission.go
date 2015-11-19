package dao

import (
	"fmt"
	"github.com/go-galaxy/galaxy/cache"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.Permission))
}

// Addmodel.Permission insert a new model.Permission into database and returns
// last inserted Id on success.
func AddPermission(m *model.Permission) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	id, err = o.Insert(m)
	if err == nil {
		m.Id = id
		cache.Permissions[id] = m
	}
	return
}

// Getmodel.PermissionById retrieves model.Permission by Id. Returns error if
// Id doesn't exist
func GetPermissionById(id int64) (v *model.Permission, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v = &model.Permission{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllmodel.Permission retrieves all model.Permission matches certain condition. Returns empty list if
// no records exist
func GetAllPermission() (ml []*model.Permission, err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	qs := o.QueryTable("permission")
	_, err = qs.All(&ml)
	return
}

type ShowPermission struct {
	Id       int64
	Name     string
	Path     string
	ParentId int
	Sub      []*ShowPermission
}

func GetFormatPermission() []*ShowPermission {
	o := orm.NewOrm()
	o.Using(model.DbName)
	var ml []*model.Permission
	qs := o.QueryTable("permission")
	_, err := qs.All(&ml)
	if err != nil {
		return nil
	}
	var List []*ShowPermission
	for _, v := range ml {
		if v.ParentId == 0 {
			List = append(List, &ShowPermission{Id: v.Id, Name: v.Name, Path: v.Path, ParentId: v.ParentId})
		}
	}

	for k, v := range List {
		for _, m := range ml {
			if v.Id == int64(m.ParentId) {
				List[k].Sub = append(List[k].Sub, &ShowPermission{Id: m.Id, Name: m.Name, Path: m.Path, ParentId: m.ParentId})
			}
		}
	}

	for _, v := range List {
		for h, l := range v.Sub {
			for _, m := range ml {
				if l.Id == int64(m.ParentId) {
					v.Sub[h].Sub = append(v.Sub[h].Sub, &ShowPermission{Id: m.Id, Name: m.Name, Path: m.Path, ParentId: m.ParentId})
				}
			}
		}
	}

	for _, v := range List {
		for h, l := range v.Sub {
			for m, n := range l.Sub {
				for _, f := range ml {
					if n.Id == int64(f.ParentId) {
						v.Sub[h].Sub[m].Sub = append(v.Sub[h].Sub[m].Sub, &ShowPermission{Id: f.Id, Name: f.Name, Path: f.Path, ParentId: f.ParentId})

					}
				}
			}
		}
	}

	return List
}

// Updatemodel.Permission updates model.Permission by Id and returns error if
// the record to be updated doesn't exist
func UpdatePermissionById(m *model.Permission) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.Permission{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
			cache.Permissions[m.Id] = m
		}
	}
	return
}

// Deletemodel.Permission deletes model.Permission by Id and returns error if
// the record to be deleted doesn't exist
func DeletePermission(id int64) (err error) {
	o := orm.NewOrm()
	o.Using(model.DbName)
	v := model.Permission{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&model.Permission{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
			delete(cache.Permissions, id)
		}
	}
	return
}
