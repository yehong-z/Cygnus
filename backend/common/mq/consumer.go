package mq

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

type Consumer interface {
	Start()
	Close()
}

type Message interface {
	Payload() []byte
	Raw() interface{}
	Context() context.Context
}

type Processor interface {
	Add(msg *sarama.ConsumerMessage)
}

type KafkaConsumer struct {
	startOnce sync.Once
	closeOnce sync.Once
	processor Processor
	config    *sarama.Config
	topic     string
	consumer  sarama.Consumer
}

func NewKafkaConsumer(processor Processor, config *sarama.Config, broker []string, topic string) *KafkaConsumer {
	if config == nil {
		config = sarama.NewConfig()
	}

	consumer, err := sarama.NewConsumer(broker, config)
	if err != nil {
		panic(err)
	}

	return &KafkaConsumer{processor: processor, config: config, consumer: consumer, topic: topic}
}

func (r *KafkaConsumer) Start() {
	r.startOnce.Do(func() {
		partitions, err := r.consumer.Partitions(r.topic)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(partitions)
		// 启动消费者协程
		for _, partition := range partitions {
			partitionConsumer, err := r.consumer.ConsumePartition(r.topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}

			go func(pc sarama.PartitionConsumer) {
				defer pc.AsyncClose()
				fmt.Println("start")
				for message := range pc.Messages() {
					fmt.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n",
						message.Partition, message.Offset, string(message.Key), string(message.Value))
					// 在这里处理接收到的消息
					r.processor.Add(message)
				}
			}(partitionConsumer)
		}
	})
}

func (r *KafkaConsumer) Close() {

}
