package main

import (
	"CoinRobot/logger"
	"CoinRobot/mq"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

var log = logger.NewLog()

var supportCoinList = []string{"#比特币/#BTC", "#以太坊/#ETH", "#莱特币/#LTC", "#柚子币/#EOS", "#比特现金/#BCH", "#瑞波币/#XRP", "#波卡币#/#DOT", "#艾达币/#ADA", "#LINK", "#币安币/#BNB", "#恒星币/#XLM", "#WBTC", "#比特币SV/#BSV", "#AAVE", "#门罗币/#XMR", "#UNI", "#SNX", "#XTZ", "#波场/#TRX", "#唯链/#VET", "#新经币/#XEM", "#ATOM", "#THETA", "#小蚁币/#NEO", "#CRO", "#OK积分/#OKB", "#DAI", "#LEO"}

func Handler(this *gin.Context) {
	body, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		log.Errorf("err: %v", err)
		this.JSON(http.StatusBadRequest, &Response{
			Code: http.StatusBadRequest,
			Data: err,
		})
		return
	}

	var message Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Errorf("err: %v", err)
		this.JSON(http.StatusBadRequest, &Response{
			Code: http.StatusBadRequest,
			Data: err,
		})
		return
	}

	switch message.PostType {
	case TypeMessage:
		switch message.MessageType {
		case MessageGroup:
			if message.SubType != MessageSubTypeNormal {
				log.Warnf("skip message_sub_type")
				this.JSON(http.StatusBadRequest, &Response{
					Code: http.StatusBadRequest,
					Data: "skip message_sub_type",
				})
				return
			}
			err = handlerGroup(&message)
			if err != nil {
				this.JSON(http.StatusInternalServerError, &Response{
					Code: http.StatusInternalServerError,
					Data: err,
				})
				return
			}
			return
		case MessagePrivate:
			fallthrough
		default:
			log.Warnf("skip message_type")
			this.JSON(http.StatusBadRequest, &Response{
				Code: http.StatusBadRequest,
				Data: "skip message_type",
			})
			return
		}
	default:
		log.Warnf("skip post_type")
	}
}

func handlerGroup(message *Message) error {
	if strings.Contains(message.RawMessage, "#比特币小助手") || strings.Contains(message.RawMessage, "#帮助") {
		message.RawMessage = "请按照如下格式在群内发言即可获得关注币种的最新消息:\n" + strings.Join(supportCoinList, "\n")
		err := send2Group(message, mq.MessageTypeNil)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		return nil
	}

	if strings.Contains(message.RawMessage, "#比特币") || strings.Contains(strings.ToLower(message.RawMessage), "#btc") {
		err := send2Group(message, mq.MessageTypeBTC)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		return nil
	}

	if strings.Contains(message.RawMessage, "#以太") || strings.Contains(strings.ToLower(message.RawMessage), "#eth") {
		err := send2Group(message, mq.MessageTypeETH)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		return nil
	}

	if strings.Contains(message.RawMessage, "#莱特") || strings.Contains(strings.ToLower(message.RawMessage), "#ltc") {
		err := send2Group(message, mq.MessageTypeLTC)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		return nil
	}

	if strings.Contains(message.RawMessage, "#柚子") || strings.Contains(strings.ToLower(message.RawMessage), "#eos") {
		err := send2Group(message, mq.MessageTypeEOS)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		return nil
	}
	return nil
}

func send2Group(message *Message, typ string) error {
	data := mq.SenderMqMsg{
		Type: mq.ChannelNameGroup,
		Data: &mq.SendMessage{
			SendURL:     "http://127.0.0.1:5799/send_group_msg",
			SendTo:      message.GroupId,
			MessageType: typ,
			Message:     message.RawMessage,
		},
	}
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 推送mq
	err = pusher.Publish(mq.TopicName, msg)
	if err != nil {
		return err
	}

	return nil
}
