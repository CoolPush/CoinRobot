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
	MessageTypeNil   = "NIL"
	MessageTypeBTC   = "BTC"
	MessageTypeETH   = "ETH"
	MessageTypeLTC   = "LTC"
	MessageTypeEOS   = "EOS"
	MessageTypeBCH   = "BCH"
	MessageTypeXRP   = "XRP"
	MessageTypeDOT   = "DOT"
	MessageTypeADA   = "ADA"
	MessageTypeLINK  = "LINK"
	MessageTypeBNB   = "BNB"
	MessageTypeXLM   = "XLM"
	MessageTypeWBTC  = "WBTC"
	MessageTypeBSV   = "BSV"
	MessageTypeAAVE  = "AAVE"
	MessageTypeXMR   = "XMR"
	MessageTypeUNI   = "UNI"
	MessageTypeSNX   = "SNX"
	MessageTypeXTZ   = "XTZ"
	MessageTypeTRX   = "TRX"
	MessageTypeVET   = "VET"
	MessageTypeXEM   = "XEM"
	MessageTypeATOM  = "ATOM"
	MessageTypeTHETA = "THETA"
	MessageTypeNEO   = "NEO"
	MessageTypeCRO   = "CRO"
	MessageTypeOKB   = "OKB"
	MessageTypeDAI   = "DAI"
	MessageTypeLEO   = "LEO"
)

const (
	PathInfoBTC   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=bitcoin&currency=cny"
	PathInfoETH   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=ethereum&currency=cny"
	PathInfoLTC   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=litecoin&currency=cny"
	PathInfoEOS   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=enterpriseOperationSystem&currency=cny"
	PathInfoBCH   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=bitcoinCash&currency=cny"
	PathInfoXRP   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=ripple&currency=cny"
	PathInfoDOT   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=polkadot&currency=cny"
	PathInfoADA   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=cardano&currency=cny"
	PathInfoLINK  = "https://www.aicoin.cn/api/coin-profile/index?coin_type=chainlink&currency=cny"
	PathInfoBNB   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=chainlink&currency=cny"
	PathInfoXLM   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=stellar&currency=cny"
	PathInfoWBTC  = "https://www.aicoin.cn/api/coin-profile/index?coin_type=wrapbtc&currency=cny"
	PathInfoBSV   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=bsv&currency=cny"
	PathInfoAAVE  = "https://www.aicoin.cn/api/coin-profile/index?coin_type=aave&currency=cny"
	PathInfoXMR   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=monero&currency=cny"
	PathInfoUNI   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=uni&currency=cny"
	PathInfoSNX   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=havven&currency=cny"
	PathInfoXTZ   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=tezos&currency=cny"
	PathInfoTRX   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=tron&currency=cny"
	PathInfoVET   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=vechain&currency=cny"
	PathInfoXEM   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=newEconomyMovement&currency=cny"
	PathInfoATOM  = "https://www.aicoin.cn/api/coin-profile/index?coin_type=atom&currency=cny"
	PathInfoTHETA = "https://www.aicoin.cn/api/coin-profile/index?coin_type=theta&currency=cny"
	PathInfoNEO   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=neo&currency=cny"
	PathInfoCRO   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=cryptocom&currency=cny"
	PathInfoOKB   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=okb&currency=cny"
	PathInfoDAI   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=dai&currency=cny"
	PathInfoLEO   = "https://www.aicoin.cn/api/coin-profile/index?coin_type=leo&currency=cny"
)

type SendMessage struct {
	SendURL     string `json:"send_url"`
	SendTo      int64  `json:"send_to"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

type SendGroupMessage struct {
	GroupId int64  `json:"group_id"`
	Message string `json:"message"`
}

type SendPrivateMessage struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
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
