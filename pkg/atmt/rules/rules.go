package rules

import (
	"strings"

	"github.com/chiyoi/trinity/pkg/atmt"
)

type Rule func(msg atmt.Message) bool

func MessageType(typ atmt.MessageType) Rule {
	return func(msg atmt.Message) bool {
		return msg.Type == typ
	}
}

func HasPrefix(prefix string) Rule {
	return func(msg atmt.Message) bool {
		return strings.HasPrefix(msg.Plaintext(), prefix)
	}
}

func Contains(kws ...string) Rule {
	return func(msg atmt.Message) bool {
		s := msg.Plaintext()
		for _, kw := range kws {
			if strings.Contains(s, kw) {
				return true
			}
		}
		return false
	}
}

func ExactlyOneOf(msgs ...string) Rule {
	m := make(map[string]bool)
	for _, msg := range msgs {
		m[msg] = true
	}
	return func(msg atmt.Message) bool {
		return m[msg.Plaintext()]
	}
}

func And(rs ...Rule) Rule {
	return func(msg atmt.Message) bool {
		for _, r := range rs {
			if !r(msg) {
				return false
			}
		}
		return true
	}
}

func Or(rs ...Rule) Rule {
	return func(msg atmt.Message) bool {
		for _, r := range rs {
			if r(msg) {
				return true
			}
		}
		return false
	}
}
