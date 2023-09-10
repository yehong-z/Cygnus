package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO: 配置写进Nacos

const MysqlConfig = "root:zyh130452@(121.36.89.81:3306)/reply?charset=utf8mb4&parseTime=True&loc=Local"

func NewMysql() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(MysqlConfig))
}
