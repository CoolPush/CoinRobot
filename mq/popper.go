package mq

import (
	"encoding/json"
	"errors"
	"github.com/guonaihong/gout"
	"github.com/nsqio/go-nsq"
	"strconv"
	"strings"
	"time"
)

type Popper struct {
	consumer *nsq.Consumer
}

func NewPopper(topicName, channelName string) (*Popper, error) {
	var consumer *nsq.Consumer
	var err error
	var config = nsq.NewConfig()
	config.LookupdPollInterval = time.Second * 30
	if consumer, err = nsq.NewConsumer(topicName, channelName, config); err != nil {
		return nil, err
	}

	consumer.SetLoggerLevel(nsq.LogLevelWarning)

	return &Popper{
		consumer: consumer,
	}, nil
}

func (pop *Popper) AddHandler() {
	pop.consumer.AddHandler(pop)
}

func (pop *Popper) ConnectToNSQLookupd() error {
	err := pop.consumer.ConnectToNSQLookupd(ConsumerAddr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (pop *Popper) HandleMessage(msg *nsq.Message) error {
	defer msg.Finish()

	var event SenderMqMsg
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	var news string

	var m SendMessage
	msgBody, err := json.Marshal(event.Data)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	err = json.Unmarshal(msgBody, &m)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	if path, exist := CoinMap[m.MessageType]; exist {
		if path != "" {
			news, err = getCoinInfo(path)
		} else {
			news = m.Message
		}
	} else {
		log.Warnf("unsupport type: %v", m.MessageType)
		return nil
	}

	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	m.Message = news

	switch event.Type {
	case ChannelNameSingle:
		err = pop.sendSingleMessage(&m)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	case ChannelNameGroup:
		err = pop.sendGroupMessage(&m)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	default:
		log.Warnf("skip event type")
	}
	return nil
}

func (pop *Popper) sendSingleMessage(data *SendMessage) error {
	if data.Message == "" {
		log.Errorf("get message empty")
		return errors.New("get message empty")
	}

	//发起推送

	var pushRet QQMessageSendResult

	err := gout.POST(data.SendURL).SetJSON(gout.H{
		"user_id": data.SendTo,
		"message": data.Message,
	}).BindJSON(&pushRet).Do()

	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	if pushRet.Retcode != 0 {
		log.Errorf("err: %+v", pushRet)
		return errors.New("推送异常")
	}

	return nil
}

func (pop *Popper) sendGroupMessage(data *SendMessage) error {
	if data.Message == "" {
		log.Errorf("get message empty")
		return errors.New("get message empty")
	}

	//发起推送
	var pushRet QQMessageSendResult

	err := gout.POST(data.SendURL).SetJSON(gout.H{
		"group_id": data.SendTo,
		"message":  data.Message,
	}).BindJSON(&pushRet).Do()

	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	if pushRet.Retcode != 0 {
		log.Errorf("push rsp: %+v", pushRet)
		return errors.New("推送异常")
	}

	return nil
}

func getCoinInfo(src string) (string, error) {
	var infoRsp CoinBase
	err := gout.GET(src).SetHeader(gout.H{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		"Referer":    "https://www.aicoin.cn/",
		"Host":       "www.aicoin.cn"}).BindJSON(&infoRsp).Do()
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	if !infoRsp.Ok {
		log.Errorf("get rsp: %+v", infoRsp)
		return "", errors.New("response err")
	}

	if infoRsp.Info == nil {
		log.Errorf("get rsp: %+v", infoRsp)
		return "", errors.New("response empty")
	}
	var info = infoRsp.Info
	content := "当前币种: " + info.Name + "\n当前美元价位: " + info.LastUsd + "$\n当前RMB价位: " + info.LastCny + "￥\n24小时涨幅: " + info.Degree24H + "%\n涨幅金额: " + info.Change24H + "￥\n多空博弈: " + getOrient(info.Orient) + "\n多空占比: " + info.Orient + "%\n市值排名: 顺" + info.Rank + "位\n当前市值: " + getSupplyValue(info.SupplyValue) + "\n标签: " + strings.Join(info.Labels, ",")
	return content, nil
}

func getOrient(orient string) string {
	if strings.Contains(orient, "-") {
		return "空方主导"
	}
	return "多方主导"
}

func getSupplyValue(val string) string {
	v, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		log.Errorf("err； %v", err)
		return "未知"
	}
	val = strconv.FormatFloat(float64(v)/100000000, 'f', 2, 64) + "亿"
	return val
}
