package rules

import (
	"github.com/chiyoi/trinity/pkg/atmt"
)

type Rule func(msg atmt.Message) bool

func Type(typ atmt.MessageType) Rule {
	return func(msg atmt.Message) bool {
		return msg.Type == typ
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
