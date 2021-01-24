package mq

import (
	"errors"
	"github.com/nsqio/go-nsq"
)

type Pusher struct {
	producer *nsq.Producer
}

func NewPusher() (*Pusher, error) {
	var producer *nsq.Producer
	var err error
	var config = nsq.NewConfig()
	if producer, err = nsq.NewProducer(ProducerAddr, config); err != nil {
		return nil, err
	}

	return &Pusher{
		producer: producer,
	}, nil
}

func (pusher *Pusher) Publish(topicName string, data []byte) error {
	if data == nil {
		log.Warnf("message is empty")
		return errors.New("message is empty")
	}

	if err := pusher.producer.Publish(topicName, data); err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
