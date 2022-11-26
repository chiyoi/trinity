package onebot

import (
	"time"
)

type EventType string

const (
	EventMessage   EventType = "message"
	EventNotice    EventType = "notice"
	EventRequest   EventType = "request"
	EventMetaEvent EventType = "meta_event"
)

type EventMessageType string

const (
	MessageGroup   EventMessageType = "group"
	MessagePrivate EventMessageType = "private"
)

type MessageSubtype string

const (
	MessageFriend    MessageSubtype = "friend"
	MessageGroupSelf MessageSubtype = "group_self"
	MessageNormal    MessageSubtype = "normal"
	MessageAnonymous MessageSubtype = "anonymous"
	MessageNotice    MessageSubtype = "notice"
)

type NoticeType string

const (
	NoticeGroupUpload   NoticeType = "group_upload"
	NoticeGroupAdmin    NoticeType = "group_admin"
	NoticeGroupDecrease NoticeType = "group_decrease"
	NoticeGroupIncrease NoticeType = "group_increase"
	NoticeGroupBan      NoticeType = "group_ban"
	NoticeFriendAdd     NoticeType = "friend_add"
	NoticeGroupRecall   NoticeType = "group_recall"
	NoticeFriendRecall  NoticeType = "friend_recall"
	NoticeGroupCard     NoticeType = "group_card"
	NoticeOfflineFile   NoticeType = "offline_file"
	NoticeClientStatus  NoticeType = "client_status"
	NoticeEssence       NoticeType = "essence"
	NoticeNotify        NoticeType = "notify"
)

type NoticeSubtype string

const (
	NoticeNotifyPoke      NoticeSubtype = "poke"
	NoticeNotifyHonor     NoticeSubtype = "honor"
	NoticeNotifyLuckyKing NoticeSubtype = "lucky_king"

	NoticeGroupAdminSet   NoticeSubtype = "set"
	NoticeGroupAdminUnset NoticeSubtype = "unset"

	NoticeGroupDecreaseLeave  NoticeSubtype = "leave"
	NoticeGroupDecreaseKick   NoticeSubtype = "kick"
	NoticeGroupDecreaseKickMe NoticeSubtype = "kick_me"

	NoticeGroupIncreaseApprove NoticeSubtype = "approve"
	NoticeGroupIncreaseInvite  NoticeSubtype = "invite"

	NoticeEventGroupBanBan     NoticeSubtype = "ban"
	NoticeEventGroupBanLiftBan NoticeSubtype = "lift_ban"

	NoticeEssenceAdd    NoticeSubtype = "add"
	NoticeEssenceDelete NoticeSubtype = "delete"
)

type HonorType string

const (
	HonorTalkative HonorType = "talkative"
	HonorPerformer HonorType = "performer"
	HonorEmotion   HonorType = "emotion"
)

type RequestType string

const (
	RequestFriend RequestType = "friend"
	RequestGroup  RequestType = "group"
)

type RequestSubtype string

const (
	RequestGroupSubtypeAdd    RequestSubtype = "add"
	RequestGroupSubtypeInvite RequestSubtype = "invite"
)

type MetaEventType string

const (
	MetaeventTypeLifecycle MetaEventType = "lifecycle"
	MetaeventTypeHeartbeat MetaEventType = "heartbeat"
)

type Event struct {
	Time     int64     `json:"time"`
	SelfId   UserId    `json:"self_id"`
	PostType EventType `json:"post_type"`

	// message event
	MessageType EventMessageType `json:"message_type"`
	SubType     MessageSubtype   `json:"sub_type"`
	MessageId   MessageId        `json:"message_id"`
	UserId      UserId           `json:"user_id"`
	Message     Message          `json:"message"`
	RawMessage  string           `json:"raw_message"`
	Font        int              `json:"font"`
	Sender      Sender           `json:"sender"`
	// primary message
	TempSource int32 `json:"temp_source"`
	// group message
	GroupId   GroupId `json:"group_id"`
	Anonymous struct {
		Id   UserId `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`

	// notice event
	NoticeType NoticeType `json:"notice_type"`

	// group upload
	// GroupId GroupId `json:"group_id"`
	// UserId  UserId  `json:"user_id"`
	File struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Size  int64  `json:"size"`
		BusId int64  `json:"busid"`
	} `json:"file"`

	// group admin
	// SubType NoticeSubtype `json:"sub_type"`
	// GroupId GroupId       `json:"group_id"`
	// UserId  UserId        `json:"user_id"`

	// group decrease
	// SubType    NoticeSubtype `json:"sub_type"`
	// GroupId    GroupId       `json:"group_id"`
	OperatorId UserId `json:"operator_id"`
	// UserId     UserId        `json:"user_id"`

	// group increase
	// SubType    NoticeSubtype `json:"sub_type"`
	// GroupId    GroupId       `json:"group_id"`
	// OperatorId UserId        `json:"operator_id"`
	// UserId     UserId        `json:"user_id"`

	// group ban
	// SubType    NoticeSubtype `json:"sub_type"`
	// GroupId    GroupId       `json:"group_id"`
	// OperatorId UserId        `json:"operator_id"`
	// UserId     UserId        `json:"user_id"`
	// Duration   int64         `json:"duration"`

	// friend add
	// UserId UserId `json:"user_id"`

	// group recall
	// GroupId    GroupId   `json:"group_id"`
	// UserId     UserId    `json:"user_id"`
	// OperatorId UserId    `json:"operator_id"`
	// MessageId  MessageId `json:"message_id"`

	// friend recall
	// UserId    UserId    `json:"user_id"`
	// MessageId MessageId `json:"message_id"`

	// notify
	// SubType NoticeSubtype `json:"sub_type"`

	// notify friend poke
	// SenderId UserId `json:"sender_id"`
	// UserId   UserId `json:"user_id"`
	// TargetId UserId `json:"target_id"`

	// notify group poke
	// GroupId  GroupId `json:"group_id"`
	// UserId   UserId  `json:"user_id"`
	// TargetId UserId  `json:"target_id"`

	// notify lucky king
	// GroupId  int64 `json:"group_id"`
	// UserId   int64 `json:"user_id"`
	TargetId int64 `json:"target_id"`

	// notify honor
	// GroupId   GroupId   `json:"group_id"`
	// UserId    UserId    `json:"user_id"`
	HonorType HonorType `json:"honor_type"`

	// group card
	// GroupId GroupId `json:"group_id"`
	// UserId  UserId  `json:"user_id"`
	CardNew string `json:"card_new"`
	CardOld string `json:"card_old"`

	// offline file
	// UserId UserId `json:"user_id"`
	// File   struct {
	//	Name string `json:"name"`
	//	Size int64  `json:"size"`
	//	Url  string `json:"url"`
	// } `json:"file"`

	// client status
	Client struct {
		AppId      int64  `json:"app_id"`
		DeviceName string `json:"device_name"`
		DeviceKind string `json:"device_kind"`
	} `json:"client"`
	Online bool `json:"online"`

	// essence
	// SubType    NoticeSubtype `json:"sub_type"`
	SenderId UserId `json:"sender_id"`
	// OperatorId UserId        `json:"operator_id"`
	// MessageId  MessageId     `json:"message_id"`

	// request
	RequestType RequestType `json:"request_type"`

	// friend
	// UserId  UserId `json:"user_id"`
	// Comment string `json:"comment"`
	// Flag    string `json:"flag"`

	// group
	// SubType RequestSubtype `json:"sub_type"`
	// GroupId GroupId        `json:"group_id"`
	// UserId  UserId         `json:"user_id"`
	Comment string `json:"comment"`
	Flag    string `json:"flag"`

	// meta event
	MetaEventType MetaEventType `json:"meta_event_type"`
	// lifecycle
	// SubType MessageSubtype `json:"sub_type"`
	// heartbeat
	Status   map[string]any `json:"status"`
	Interval time.Duration  `json:"interval"`
}

type (
	UserId    = int64
	GroupId   = int64
	MessageId = int32
)

type Sex string

const (
	SexMale    Sex = "male"
	SexFemale  Sex = "female"
	SexUnknown Sex = "unknown"
)

type Sender struct {
	UserId   UserId `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      Sex    `json:"sex"`
	Age      int32  `json:"age"`
}

type MessageType string

const (
	MessageText      MessageType = "text"
	MessageFace      MessageType = "face"
	MessageRecord    MessageType = "record"
	MessageVideo     MessageType = "video"
	MessageAt        MessageType = "at"
	MessageMusic     MessageType = "music"
	MessageImage     MessageType = "image"
	MessageReply     MessageType = "reply"
	MessageRedbag    MessageType = "redbag"
	MessagePoke      MessageType = "poke"
	MessageGift      MessageType = "gift"
	MessageForward   MessageType = "forward"
	MessageNode      MessageType = "node"
	MessageXML       MessageType = "xml"
	MessageJSON      MessageType = "json"
	MessageCardimage MessageType = "cardimage"
	MessageTTS       MessageType = "tts"
)

type Message []MessageSegment

type MessageSegment struct {
	Type MessageType       `json:"type"`
	Data map[string]string `json:"data"`
}

func (msg Message) Append(seg MessageSegment) Message { return append(msg, seg) }
func (msg Message) Extend(msg1 Message) Message       { return append(msg, msg1...) }

func (seg MessageSegment) Chain(seg1 MessageSegment) Message { return append(Message{}, seg, seg1) }
