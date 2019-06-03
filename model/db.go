package model

import (
	"fmt"
	"sku-manage/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB

func DBInit() {

	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8",
		config.SrvConfig.Mysql.User,
		config.SrvConfig.Mysql.Password,
		config.SrvConfig.Mysql.Host,
		config.SrvConfig.Mysql.Port,
		config.SrvConfig.Mysql.DbName,
	))
	if err != nil {
		logrus.Debugf("error %s", err.Error())
		panic("failed to connect database")
	}

	db.Set("gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		&User{},
		&Sku{},
		&Group{},
		&Company{},
		&Role{},
		&Department{},
		&Purchase{},
		&URL{},
	)

	db.Model(&Department{}).AddUniqueIndex("idx_user_name_age", "name", "size")
	db.LogMode(true)

}
