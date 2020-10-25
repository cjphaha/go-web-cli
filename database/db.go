package database

import (
	"database/sql"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
)

var (
	cfg, _ = goconfig.LoadConfigFile("config.ini");
	name, _      = cfg.GetValue("mysql", "username");
	password,_ = cfg.GetValue("mysql", "password");
	url,_ = cfg.GetValue("mysql", "url");
	Db, Err = sql.Open("mysql", name+":"+password+"@tcp"+url+"?parseTime=true");
	DB *gorm.DB
)

func InintDatabase(){
	db, err := gorm.Open("mysql",name+":"+password+"@tcp"+url+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Println(err)
	}
	//SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.DB().SetMaxOpenConns(10)   //设置数据库连接池最大连接数
	db.DB().SetMaxIdleConns(4)   //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	DB = db
}

func GetDB() *gorm.DB{
	return DB
}
