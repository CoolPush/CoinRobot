package main

import (
	"CoinRobot/logger"
	"CoinRobot/mq"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"net/http"
	"strings"
)

var log = logger.NewLog()

var supportCoinList = []string{"#比特币/#BTC", "#以太坊/#ETH", "#莱特币/#LTC", "#柚子币/#EOS", "#比特现金/#BCH", "#瑞波币/#XRP", "#波卡币#/#DOT", "#LINK", "#比特币SV/#BSV", "#门罗币/#XMR", "#UNI", "#波场/#TRX", "#THETA"}

func HandlerCoin(this *gin.Context) {
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

	msg, ok := message.Message.Message.(string)
	if !ok {
		return nil
	}

	if strings.Contains(msg, "#助手") || strings.Contains(msg, "#帮助") {
		message.RawMessage = "请按照如下格式发言即可获得关注币种的最新消息:\n" + strings.Join(supportCoinList, "\n")
		coinType = mq.MessageTypeNil
	} else if strings.Contains(msg, "#比特") || strings.Contains(strings.ToLower(msg), "#btc") {
		coinType = mq.MessageTypeBTC
	} else if strings.Contains(msg, "#以太") || strings.Contains(strings.ToLower(msg), "#eth") {
		coinType = mq.MessageTypeETH
	} else if strings.Contains(msg, "#莱特") || strings.Contains(strings.ToLower(msg), "#ltc") {
		coinType = mq.MessageTypeLTC
	} else if strings.Contains(msg, "#柚子") || strings.Contains(strings.ToLower(msg), "#eos") {
		coinType = mq.MessageTypeEOS
	} else if strings.Contains(msg, "#比特现金") || strings.Contains(strings.ToLower(msg), "#bch") {
		coinType = mq.MessageTypeBCH
	} else if strings.Contains(msg, "#瑞波") || strings.Contains(strings.ToLower(msg), "#xrp") {
		coinType = mq.MessageTypeXRP
	} else if strings.Contains(msg, "#波卡") || strings.Contains(strings.ToLower(msg), "#dot") {
		coinType = mq.MessageTypeDOT
	} else if strings.Contains(strings.ToLower(msg), "#link") {
		coinType = mq.MessageTypeLINK
	} else if strings.Contains(msg, "#SV") || strings.Contains(strings.ToLower(msg), "#bsv") {
		coinType = mq.MessageTypeBSV
	} else if strings.Contains(msg, "#门罗") || strings.Contains(strings.ToLower(msg), "#xmr") {
		coinType = mq.MessageTypeXMR
	} else if strings.Contains(msg, "#UNI") || strings.Contains(strings.ToLower(msg), "#uni") {
		coinType = mq.MessageTypeUNI
	} else if strings.Contains(msg, "#波场") || strings.Contains(strings.ToLower(msg), "#trx") {
		coinType = mq.MessageTypeTRX
	} else if strings.Contains(strings.ToLower(msg), "#theta") {
		coinType = mq.MessageTypeTHETA
	} else {
		return nil
	}

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

// -= HandlerLSP

func HandlerLSP(this *gin.Context) {
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
			err = handlerSendLSP(&message, MessageGroup)
		case MessagePrivate:
			err = handlerSendLSP(&message, MessagePrivate)
		default:
			log.Warnf("skip message_type")
			this.JSON(http.StatusBadRequest, &Response{
				Code: http.StatusBadRequest,
				Data: "skip message_type",
			})
			return
		}
	default:
		break
	}

	if err != nil {
		this.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Data: err,
		})
		return
	}
}

func handlerSendLSP(message *PostMessage, sendType string) error {
	msg, ok := message.Message.Message.(string)
	if !ok {
		return nil
	}

	if strings.Contains(msg, "开车") {
		message.RawMessage = "### 开车 ###\n输入 车来 即可完成开车\n资源来自妹子图，不保证质量"
	} else if strings.Contains(msg, "车来") {
		img, err := getPIC()
		if err != nil {
			return err
		}
		message.RawMessage = fmt.Sprintf("[CQ:image,file=%s]", img)
	} else {
		return nil
	}

	switch sendType {
	case MessageGroup:
		err := send2Group(message, mq.MessageTypeNil)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	case MessagePrivate:
		err := send2Single(message, mq.MessageTypeNil)
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
	}

	return nil
}

func getPIC() (string, error) {
	var resp GetMzPicResponse
	err := gout.GET("https://api.func.ws/api/img/mz?format=json").
		BindJSON(&resp).Do()
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	return resp.Data.Img, nil
}
