package main

import (
	"CoinRobot/mq"
	"os"
)

var pusher *mq.Pusher

func init() {
	var err error
	pusher, err = mq.NewPusher()
	if err != nil {
		log.Errorf("err: %v", err)
		panic(err)
	}

	go initConsumer()
}

func initConsumer() {
	{
		popper, err := mq.NewPopper(mq.TopicName, mq.ChannelNameGroup)
		if err != nil {
			log.Fatal(err)
		}
		popper.AddHandler()
		err = popper.ConnectToNSQLookupd()
		if err != nil {
			log.Errorf("err: %v", err)
			os.Exit(1)
		}
	}

	{
		popper2, err := mq.NewPopper(mq.TopicName, mq.ChannelNameGroup)
		if err != nil {
			log.Fatal(err)
		}
		popper2.AddHandler()
		err = popper2.ConnectToNSQLookupd()
		if err != nil {
			log.Errorf("err: %v", err)
			os.Exit(1)
		}
	}

	{
		popper3, err := mq.NewPopper(mq.TopicName, mq.ChannelNameGroup)
		if err != nil {
			log.Fatal(err)
		}
		popper3.AddHandler()
		err = popper3.ConnectToNSQLookupd()
		if err != nil {
			log.Errorf("err: %v", err)
			os.Exit(1)
		}
	}
	<-make(chan bool)
}
