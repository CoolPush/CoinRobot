package mq

import (
	"CoinRobot/logger"
)

var log = logger.NewLog()

const (
	ProducerAddr = "127.0.0.1:4150"
	ConsumerAddr = "127.0.0.1:4161"
)

const (
	TopicName         = "coin_robot"
	ChannelNameSingle = "single"
	ChannelNameGroup  = "group"
)

type SenderMqMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"message"`
}

const (
	MessageTypeBTC = "BTC"
	MessageTypeETH = "ETH"
	MessageTypeLTC = "LTC"
	MessageTypeEOS = "EOS"
)

type SendMessage struct {
	SendURL     string `json:"send_url"`
	SendTo      int64  `json:"send_to"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

type CoinInfo struct {
	Rank        string   `json:"rank"`
	Name        string   `json:"name"`
	Logo        string   `json:"logo"`
	LastCny     string   `json:"last_cny"`
	LastUsd     string   `json:"last_usd"`
	Degree24H   string   `json:"degree_24h"`
	Change24H   string   `json:"change_24h"`
	Vol24H      string   `json:"vol24h"`
	Trade24H    string   `json:"trade24h"`
	Orient      string   `json:"orient"`
	SupplyValue string   `json:"supply_value"`
	UpPercent   string   `json:"up_percent"`
	Labels      []string `json:"labels"`
}

type CoinBase struct {
	Ok   bool      `json:"ok"`
	Info *CoinInfo `json:"global"`
}
