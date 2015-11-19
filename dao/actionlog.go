package dao

import (
	"fmt"
	"github.com/go-galaxy/galaxy/model"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.ActionLog))
}

// Addmodel.ActionLog insert a new model.ActionLog into database
//and returns error.
func AddActionLog(m *model.ActionLog) error {
	o := orm.NewOrm()
	o.Using(model.DbName)
	fmt.Println(m)
	id, err := o.Insert(m)
	fmt.Println(id)
	fmt.Println(err)
	return err
}
