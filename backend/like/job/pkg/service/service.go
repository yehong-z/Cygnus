package service

import "github.com/yehong-z/Cygnus/common/mq"

func InitLikeProcessor() {
	k := mq.NewKafkaConsumer(nil, nil, []string{"121.36.89.81:9092"}, "test")
	k.Start()
}
