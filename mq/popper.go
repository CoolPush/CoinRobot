package mq

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/guonaihong/gout"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
	"net/http"
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
	switch event.Type {
	case ChannelNameSingle:
		var msg SendMessage
		msgBody, err := json.Marshal(event.Data)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		err = json.Unmarshal(msgBody, &msg)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		err = pop.sendSingleMessage(&msg)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	case ChannelNameGroup:
		var msg SendMessage
		msgBody, err := json.Marshal(event.Data)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		err = json.Unmarshal(msgBody, &msg)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}

		switch msg.MessageType {
		case MessageTypeBTC:
			news, err := getBTCInfo()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			msg.Message = news
		case MessageTypeETH:
			news, err := getETHInfo()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			msg.Message = news
		case MessageTypeLTC:
			news, err := getLTCInfo()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			msg.Message = news
		case MessageTypeEOS:
			news, err := getEOSInfo()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			msg.Message = news
		case MessageTypeNil:
		default:
			log.Warnf("unsupport type")
			return nil
		}

		err = pop.sendGroupMessage(&msg)
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
	//发起推送
	var pushRet = &struct {
		RetCode int64  `json:"retcode"`
		Status  string `json:"status"`
	}{}

	resp, err := http.Post(data.SendURL, "application/x-www-form-urlencoded", strings.NewReader("user_id="+strconv.FormatInt(data.SendTo, 10)+"&message="+data.Message))
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	err = json.Unmarshal(content, pushRet)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	if pushRet.RetCode != 0 {
		return errors.New("推送异常")

	}

	return nil
}

func (pop *Popper) sendGroupMessage(data *SendMessage) error {
	log.Infof("get sendGroupMessage data: %+v", data)
	if data.Message == "" {
		log.Errorf("get message empty")
		return errors.New("get message empty")
	}

	var body = SendGroupMessage{
		GroupId: data.SendTo,
		Message: data.Message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	log.Infof("json body: %v", body)

	reqBody := bytes.NewBuffer(jsonBody)

	//发起推送
	var pushRet = &struct {
		RetCode int64  `json:"retcode"`
		Status  string `json:"status"`
	}{}
	resp, err := http.Post(data.SendURL, "application/json", reqBody)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	err = json.Unmarshal(content, pushRet)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	if pushRet.RetCode != 0 {
		log.Errorf("push rsp: %+v", pushRet)
		return errors.New("推送异常")
	}

	return nil
}

func getBTCInfo() (string, error) {
	var infoUrl = "https://www.aicoin.cn/api/coin-profile/index?coin_type=bitcoin&currency=cny"
	var infoRsp CoinBase
	err := gout.GET(infoUrl).SetHeader(gout.H{
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
	log.Infof("rsp: %v", info)
	content := "当前币种: " + info.Name + "\n当前美元价位: " + info.LastUsd + "$\n当前RMB价位: " + info.LastCny + "￥\n24小时涨幅: " + info.Degree24H + "%\n涨幅金额: " + info.Change24H + "￥\n多空博弈: " + getOrient(info.Orient) + "多空占比: " + info.Orient + "%\n市值排名: 顺" + info.Rank + "位\n当前市值: " + getSupplyValue(info.SupplyValue) + "\n标签: " + strings.Join(info.Labels, ",")
	return content, nil
}

func getETHInfo() (string, error) {
	var infoUrl = "https://www.aicoin.cn/api/coin-profile/index?coin_type=ethereum&currency=cny"
	var infoRsp CoinBase
	err := gout.GET(infoUrl).SetHeader(gout.H{
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
	log.Infof("rsp: %v", info)
	content := "当前币种: " + info.Name + "\n当前美元价位: " + info.LastUsd + "$\n当前RMB价位: " + info.LastCny + "￥\n24小时涨幅: " + info.Degree24H + "%\n涨幅金额: " + info.Change24H + "￥\n多空博弈: " + getOrient(info.Orient) + "多空占比: " + info.Orient + "%\n市值排名: 顺" + info.Rank + "位\n当前市值: " + getSupplyValue(info.SupplyValue) + "\n标签: " + strings.Join(info.Labels, ",")
	return content, nil
}

func getLTCInfo() (string, error) {
	var infoUrl = "https://www.aicoin.cn/api/coin-profile/index?coin_type=bitcoin&currency=cny"
	var infoRsp CoinBase
	err := gout.GET(infoUrl).SetHeader(gout.H{
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
	log.Infof("rsp: %v", info)
	content := "当前币种: " + info.Name + "\n当前美元价位: " + info.LastUsd + "$\n当前RMB价位: " + info.LastCny + "￥\n24小时涨幅: " + info.Degree24H + "%\n涨幅金额: " + info.Change24H + "￥\n多空博弈: " + getOrient(info.Orient) + "多空占比: " + info.Orient + "%\n市值排名: 顺" + info.Rank + "位\n当前市值: " + getSupplyValue(info.SupplyValue) + "\n标签: " + strings.Join(info.Labels, ",")
	return content, nil
}

func getEOSInfo() (string, error) {
	var infoUrl = "https://www.aicoin.cn/api/coin-profile/index?coin_type=enterpriseOperationSystem&currency=cny"
	var infoRsp CoinBase
	err := gout.GET(infoUrl).SetHeader(gout.H{
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
	log.Infof("rsp: %v", info)
	content := "当前币种: " + info.Name + "\n当前美元价位: " + info.LastUsd + "$\n当前RMB价位: " + info.LastCny + "￥\n24小时涨幅: " + info.Degree24H + "%\n涨幅金额: " + info.Change24H + "￥\n多空博弈: " + getOrient(info.Orient) + "多空占比: " + info.Orient + "%\n市值排名: 顺" + info.Rank + "位\n当前市值: " + getSupplyValue(info.SupplyValue) + "\n标签: " + strings.Join(info.Labels, ",")
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
