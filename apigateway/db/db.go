package db

import (
	"apigateway/cfg"

	"github.com/qq51529210/util"
	"gorm.io/gorm"
)

var (
	_db *gorm.DB
)

// 用于取地址
var (
	True  int8 = 1
	False int8 = 0
)

// Init 初始化数据库
func Init() error {
	var err error
	// 初始化
	dbCfg := util.NewGORMConfig()
	dbCfg.Logger = &util.GORMLog{}
	_db, err = util.InitGORM(cfg.Cfg.DB.URL, dbCfg)
	if err != nil {
		return err
	}
	// 数据表
	err = initTable()
	if err != nil {
		return err
	}
	//
	return nil
}

// initTable 初始化数据表
func initTable() error {
	_db.AutoMigrate(
		new(Service),
		new(Server),
		new(Admin),
	)
	//
	return nil
}
