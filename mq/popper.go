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
	var news string
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
			news, err = getCoinInfo(PathInfoBTC)
		case MessageTypeETH:
			news, err = getCoinInfo(PathInfoETH)
		case MessageTypeLTC:
			news, err = getCoinInfo(PathInfoLTC)
		case MessageTypeBCH:
			news, err = getCoinInfo(PathInfoBCH)
		case MessageTypeXRP:
			news, err = getCoinInfo(PathInfoXRP)
		case MessageTypeDOT:
			news, err = getCoinInfo(PathInfoDOT)
		case MessageTypeADA:
			news, err = getCoinInfo(PathInfoADA)
		case MessageTypeLINK:
			news, err = getCoinInfo(PathInfoLINK)
		case MessageTypeBNB:
			news, err = getCoinInfo(PathInfoBNB)
		case MessageTypeXLM:
			news, err = getCoinInfo(PathInfoXLM)
		case MessageTypeWBTC:
			news, err = getCoinInfo(PathInfoWBTC)
		case MessageTypeBSV:
			news, err = getCoinInfo(PathInfoBSV)
		case MessageTypeEOS:
			news, err = getCoinInfo(PathInfoEOS)
		case MessageTypeAAVE:
			news, err = getCoinInfo(PathInfoAAVE)
		case MessageTypeXMR:
			news, err = getCoinInfo(PathInfoXMR)
		case MessageTypeUNI:
			news, err = getCoinInfo(PathInfoUNI)
		case MessageTypeSNX:
			news, err = getCoinInfo(PathInfoSNX)
		case MessageTypeXTZ:
			news, err = getCoinInfo(PathInfoXTZ)
		case MessageTypeTRX:
			news, err = getCoinInfo(PathInfoTRX)
		case MessageTypeVET:
			news, err = getCoinInfo(PathInfoVET)
		case MessageTypeXEM:
			news, err = getCoinInfo(PathInfoXEM)
		case MessageTypeATOM:
			news, err = getCoinInfo(PathInfoATOM)
		case MessageTypeTHETA:
			news, err = getCoinInfo(PathInfoTHETA)
		case MessageTypeNEO:
			news, err = getCoinInfo(PathInfoNEO)
		case MessageTypeCRO:
			news, err = getCoinInfo(PathInfoCRO)
		case MessageTypeOKB:
			news, err = getCoinInfo(PathInfoOKB)
		case MessageTypeDAI:
			news, err = getCoinInfo(PathInfoDAI)
		case MessageTypeLEO:
			news, err = getCoinInfo(PathInfoLEO)
		case MessageTypeNil:
			news = msg.Message
		default:
			log.Warnf("unsupport type")
			return nil
		}
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		msg.Message = news

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
