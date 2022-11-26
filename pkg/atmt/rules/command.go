package rules

import (
	"strings"

	"github.com/chiyoi/trinity/pkg/atmt"
)

func Command(cmd string, commandPrefix []string) Rule {
	return func(ev atmt.Event) bool {
		s := ev.Message.Plaintext()
		for _, p := range commandPrefix {
			if strings.HasPrefix(s, p) && s[len(p):] == cmd {
				return true
			}
		}
		return false
	}
}

func Keywords(kws ...string) Rule {
	return func(ev atmt.Event) bool {
		s := ev.Message.Plaintext()
		for _, kw := range kws {
			if strings.Contains(s, kw) {
				return true
			}
		}
		return false
	}
}

func ExactMessageOneOf(msgs ...string) Rule {
	m := make(map[string]bool)
	for _, msg := range msgs {
		m[msg] = true
	}
	return func(ev atmt.Event) bool {
		return m[ev.Message.Plaintext()]
	}
}
