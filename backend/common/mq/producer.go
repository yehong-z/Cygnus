package mq

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

type Producer interface {
	Send(c context.Context, k string, v []byte) (err error)
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) Send(c context.Context, k string, v []byte) (err error) {
	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Key:   sarama.StringEncoder(k),
		Value: sarama.StringEncoder(v),
	}

	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	return nil
}

func (kp *KafkaProducer) Close() error {
	if kp.producer != nil {
		return kp.producer.Close()
	}
	return nil
}
