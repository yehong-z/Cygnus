package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/yehong-z/Cygnus/common/mq"
	"github.com/yehong-z/Cygnus/common/snow"
	"github.com/yehong-z/Cygnus/like/common/message"
	"github.com/yehong-z/Cygnus/like/job/pkg/dao"
	"gorm.io/gorm"
)

type CountProcessor struct {
	dao dao.Dao
}

func NewCountProcessor(dao dao.Dao) *CountProcessor {
	return &CountProcessor{dao: dao}
}

func (c *CountProcessor) Add(msg *sarama.ConsumerMessage) {
	cnt := message.CountMessage{}
	err := json.Unmarshal(msg.Value, &cnt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.TODO()
	count, err := c.dao.GetCount(ctx, cnt.ObjectId)
	if err != nil {
		return
	}

	if count == nil {
		err = c.dao.CreateCount(ctx, cnt.ObjectId, cnt.Like, cnt.Dislike)
		if err != nil {
			return
		}
	} else {
		err = c.dao.UpdateCount(ctx, cnt.ObjectId, cnt.Like, cnt.Dislike)
		if err != nil {
			return
		}
	}
}

func InitLikeProcessor(mysql *gorm.DB) {
	s := snow.NewSnow(1)
	d := dao.NewDao(mysql, s)
	p := NewCountProcessor(d)
	k := mq.NewKafkaConsumer(p, nil, []string{"121.36.89.81:9092"}, "like")
	fmt.Println("comsumer start")
	k.Start()
}
