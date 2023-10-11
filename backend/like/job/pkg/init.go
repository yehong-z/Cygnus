package pkg

import (
	"github.com/yehong-z/Cygnus/like/job/pkg/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const MysqlConfig = "root:zyh130452@(121.36.89.81:3306)/like?charset=utf8mb4&parseTime=True&loc=Local"

func NewMysql() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(MysqlConfig))
}

func Init() {
	m, err := NewMysql()
	if err != nil {
		panic(err)
	}

	service.InitLikeProcessor(m)
}
