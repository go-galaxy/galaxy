package service

import (
	"github.com/go-galaxy/galaxy/dao"
	"github.com/go-galaxy/galaxy/model"
)

func AddLog(mylog *model.ActionLog) (err error) {
	err = dao.AddActionLog(mylog)
	return
}
