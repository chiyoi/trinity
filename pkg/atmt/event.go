package atmt

import (
	"time"

	"github.com/chiyoi/trinity/pkg/atmt/message"
)

type Event struct {
	Time      time.Time
	User      string
	MessageId string
	Message   message.Message
}
