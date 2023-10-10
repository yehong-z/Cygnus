package mq

import (
	"context"
	"fmt"
	"testing"
)

func TestProducer(t *testing.T) {
	producer, err := NewKafkaProducer([]string{"121.36.89.81:9092"}, "test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()
	err = producer.Send(ctx, "aa", "1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestConsumer(t *testing.T) {
	c := NewKafkaConsumer(nil, nil, []string{"121.36.89.81:9092"}, "test")
	c.Start()
	select {}
}
