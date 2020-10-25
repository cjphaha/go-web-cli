package models

import "github.com/jinzhu/gorm"

type User struct {
	Name string
	Password string
	gorm.Model
}
//为数据库建立别名
func (v User) TableName() string {
	return "commonuser"
}