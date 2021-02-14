package main

const (
	SendGroupMsg     = "http://127.0.0.1:5799/send_group_msg"
	SendSingleMsg    = "http://127.0.0.1:5799/send_private_msg"
	ApproveFriendAdd = "http://127.0.0.1:5799/set_friend_add_request"
	ApproveGroupAdd  = "http://127.0.0.1:5799/set_group_add_request"
)

const (
	TypeMessage   = "message"
	TypeNotice    = "notice"
	TypeRequest   = "request"
	TypeMetaEvent = "meta_event"
)

const (
	MessagePrivate = "private"
	MessageGroup   = "group"

	MessageSubTypeNormal    = "normal"
	MessageSubTypeAnonymous = "anonymous"
	MessageSubTypeNotice    = "notice"
)

const (
	RequestFriendAdd = "friend_add_request"
	RequestGroupAdd  = "group_add_request"

	RequestTypeGroup  = "group"
	RequestTypeFriend = "friend"
)

type PostMessage struct {
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`
	SubType  string `json:"sub_type"`
	UserId   int64  `json:"user_id"`
	GroupId  int64  `json:"group_id"`
	Message
	Request
}

type Message struct {
	MessageType string      `json:"message_type"`
	MessageId   int32       `json:"message_id"`
	Anonymous   Anonymous   `json:"anonymous"`
	Message     interface{} `json:"message"`
	RawMessage  string      `json:"raw_message"`
	Font        int32       `json:"font"`
	Sender      Sender      `json:"sender"`
}

type Request struct {
	RequestType string `json:"request_type"`
	Comment     string `json:"comment"`
	Flag        string `json:"flag"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
	Card     string `json:"card"`
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`
	Area     string `json:"area"`
	Level    string `json:"level"`
	Role     string `json:"role"`
	Title    string `json:"title"`
}

type Anonymous struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// Response 是交付层的基本回应
type Response struct {
	Code    int         `json:"code"`    //请求状态代码
	Message interface{} `json:"message"` //请求结果提示
	Data    interface{} `json:"data"`    //请求结果与错误原因
}
