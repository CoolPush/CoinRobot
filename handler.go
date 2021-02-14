package main

import (
	"CoinRobot/logger"
	"CoinRobot/mq"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"net/http"
	"strings"
)

var log = logger.NewLog()

var supportCoinList = []string{"#比特币/#BTC", "#以太坊/#ETH", "#莱特币/#LTC", "#柚子币/#EOS", "#比特现金/#BCH", "#瑞波币/#XRP", "#波卡币#/#DOT", "#LINK", "#比特币SV/#BSV", "#门罗币/#XMR", "#UNI", "#波场/#TRX", "#THETA"}

func Handler(this *gin.Context) {
	var message PostMessage
	err := this.ShouldBind(&message)
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
			err = handlerSend(&message, MessageGroup)
		case MessagePrivate:
			err = handlerSend(&message, MessagePrivate)
		default:
			log.Warnf("skip message_type")
			this.JSON(http.StatusBadRequest, &Response{
				Code: http.StatusBadRequest,
				Data: "skip message_type",
			})
			return
		}
	case TypeRequest:
		switch message.RequestType {
		case RequestTypeFriend:
			err = gout.POST(ApproveFriendAdd).SetJSON(gout.H{
				"flag":    message.Request.Flag,
				"approve": true,
			}).Do()
		case RequestTypeGroup:
			switch message.SubType {
			case "add":
				err = gout.POST(ApproveGroupAdd).SetJSON(gout.H{
					"flag":     message.Request.Flag,
					"approve":  true,
					"sub_type": "add",
					"type":     "add",
				}).Do()
			case "invite":
				err = gout.POST(ApproveGroupAdd).SetJSON(gout.H{
					"flag":     message.Request.Flag,
					"approve":  true,
					"sub_type": "invite",
					"type":     "invite",
				}).Do()
			default:
				log.Warnf("skip request_sub_type")
				this.JSON(http.StatusBadRequest, &Response{
					Code: http.StatusBadRequest,
					Data: "skip request_sub_type",
				})
				return
			}
		}
		if err != nil {
			log.Errorf("err: %v", err)
			this.JSON(http.StatusBadRequest, &Response{
				Code:    http.StatusBadRequest,
				Message: "同意请求失败",
				Data:    err,
			})
			return
		}
	case TypeNotice:
		fallthrough
	case TypeMetaEvent:
		break
	default:
		log.Warnf("skip post_type: %v", message.PostType)
	}

	if err != nil {
		this.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Data: err,
		})
		return
	}
}

func handlerSend(message *PostMessage, sendType string) error {
	var coinType string
	if strings.Contains(message.RawMessage, "#助手") || strings.Contains(message.RawMessage, "#帮助") {
		message.RawMessage = "请按照如下格式在群内发言即可获得关注币种的最新消息:\n" + strings.Join(supportCoinList, "\n")
		coinType = mq.MessageTypeNil
	}

	if strings.Contains(message.RawMessage, "#比特") || strings.Contains(strings.ToLower(message.RawMessage), "#btc") {
		coinType = mq.MessageTypeBTC
	}

	if strings.Contains(message.RawMessage, "#以太") || strings.Contains(strings.ToLower(message.RawMessage), "#eth") {
		coinType = mq.MessageTypeETH
	}

	if strings.Contains(message.RawMessage, "#莱特") || strings.Contains(strings.ToLower(message.RawMessage), "#ltc") {
		coinType = mq.MessageTypeLTC
	}

	if strings.Contains(message.RawMessage, "#柚子") || strings.Contains(strings.ToLower(message.RawMessage), "#eos") {
		coinType = mq.MessageTypeEOS
	}

	if strings.Contains(message.RawMessage, "#比特现金") || strings.Contains(strings.ToLower(message.RawMessage), "#bch") {
		coinType = mq.MessageTypeBCH
	}

	if strings.Contains(message.RawMessage, "#瑞波") || strings.Contains(strings.ToLower(message.RawMessage), "#xrp") {
		coinType = mq.MessageTypeXRP
	}

	if strings.Contains(message.RawMessage, "#波卡") || strings.Contains(strings.ToLower(message.RawMessage), "#dot") {
		coinType = mq.MessageTypeDOT
	}

	if strings.Contains(strings.ToLower(message.RawMessage), "#link") {
		coinType = mq.MessageTypeLINK
	}

	if strings.Contains(message.RawMessage, "#SV") || strings.Contains(strings.ToLower(message.RawMessage), "#bsv") {
		coinType = mq.MessageTypeBSV
	}

	if strings.Contains(message.RawMessage, "#门罗") || strings.Contains(strings.ToLower(message.RawMessage), "#xmr") {
		coinType = mq.MessageTypeXMR
	}

	if strings.Contains(message.RawMessage, "#UNI") || strings.Contains(strings.ToLower(message.RawMessage), "#uni") {
		coinType = mq.MessageTypeUNI
	}

	if strings.Contains(message.RawMessage, "#TRX") || strings.Contains(strings.ToLower(message.RawMessage), "#trx") {
		coinType = mq.MessageTypeTRX
	}

	if strings.Contains(message.RawMessage, "#THETA") || strings.Contains(strings.ToLower(message.RawMessage), "#theta") {
		coinType = mq.MessageTypeTHETA
	}

	log.Infof("get coin_type: %+v, get message: %+v", coinType, message)
	switch sendType {
	case MessageGroup:
		err := send2Group(message, coinType)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	case MessagePrivate:
		err := send2Single(message, coinType)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	}

	return nil
}

func send2Group(message *PostMessage, typ string) error {
	data := mq.SenderMqMsg{
		Type: mq.ChannelNameGroup,
		Data: &mq.SendMessage{
			SendURL:     SendGroupMsg,
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

func send2Single(message *PostMessage, typ string) error {
	data := mq.SenderMqMsg{
		Type: mq.ChannelNameSingle,
		Data: &mq.SendMessage{
			SendURL:     SendSingleMsg,
			SendTo:      message.UserId,
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
