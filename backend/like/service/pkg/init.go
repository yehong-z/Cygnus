package pkg

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/yehong-z/Cygnus/common/mq"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao"
	"github.com/yehong-z/Cygnus/like/service/pkg/service"
)

func Init() {
	mysql, err := dao.NewMysql()
	if err != nil {
		panic(err)
	}

	producer, err := mq.NewKafkaProducer([]string{"121.36.89.81:9092"}, "like")
	if err != nil {
		panic(err)
	}

	d := dao.NewDao(mysql, producer)

	config.SetProviderService(service.NewThumbupServerImpl(d))
}
