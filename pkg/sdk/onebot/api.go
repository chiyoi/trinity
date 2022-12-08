package onebot

type Request struct {
	Action Action `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo"`
}
type Action string

const (
	ActionSendPrivateMsg        Action = "send_private_msg"
	ActionSendGroupMsg          Action = "send_group_msg"
	ActionSendForwardMsg        Action = "send_group_forward_msg"
	ActionSendMsg               Action = "send_msg"
	ActionSendPrivateForwardMsg Action = "send_private_forward_msg"

	ActionDeleteMsg Action = "delete_msg"

	ActionGetMsg        Action = "get_msg"
	ActionGetForwardMsg Action = "get_forward_msg"

	ActionSetFriendAddRequest Action = "set_friend_add_request"
	ActionSetGroupAddRequest  Action = "set_group_add_request"

	ActionGetLoginInfo Action = "get_login_info"
)

type ReqParamsSendMsg struct {
	MessageType EventMessageType `json:"message_type"`
	UserId      UserId           `json:"user_id"`
	GroupId     GroupId          `json:"group_id"`
	Message     Message          `json:"message"`
	AutoEscape  bool             `json:"auto_escape"`
}

type Status string

const (
	StatusOK     Status = "ok"
	StatusAsync  Status = "async"
	StatusFailed Status = "failed"
)

type Retcode int

const (
	RetcodeOK    Retcode = 0
	RetcodeAsync Retcode = 1
	// other retcode: failed
)

type Response[Data RespData] struct {
	Status  Status  `json:"status"`
	Retcode Retcode `json:"retcode"`
	Msg     string  `json:"msg"`
	Wording string  `json:"wording"`
	Data    Data    `json:"data"`
	Echo    string  `json:"echo"`
}

type RespData interface {
	RespDataSendMsg
}

type RespDataSendMsg struct {
	MessageId MessageId `json:"message_id"`
}
