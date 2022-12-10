package rules

import (
	"strings"

	"github.com/chiyoi/trinity/pkg/atmt"
)

func HasPrefix(prefix string) Rule {
	return And(
		MessageType(atmt.MessagePush),
		func(msg atmt.Message) bool {
			return strings.HasPrefix(msg.Content.Plaintext(), prefix)
		},
	)
}

func Contains(kws ...string) Rule {
	return And(
		MessageType(atmt.MessagePush),
		func(msg atmt.Message) bool {
			s := msg.Content.Plaintext()
			for _, kw := range kws {
				if strings.Contains(s, kw) {
					return true
				}
			}
			return false
		},
	)
}

func ExactlyOneOf(msgs ...string) Rule {
	m := make(map[string]bool)
	for _, msg := range msgs {
		m[msg] = true
	}
	return And(
		MessageType(atmt.MessagePush),
		func(msg atmt.Message) bool {
			return m[msg.Content.Plaintext()]
		},
	)
}
