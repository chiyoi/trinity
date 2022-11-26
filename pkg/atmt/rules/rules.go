package rules

import (
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
)

type Rule func(ev atmt2.Event) bool

func (r Rule) Match(ev atmt2.Event) bool { return r(ev) }
func (r Rule) Priority() int             { return 0 }

var _ atmt2.Matcher = (Rule)(nil)

func And(rs ...Rule) Rule {
	return func(ev atmt2.Event) bool {
		for _, r := range rs {
			if !r(ev) {
				return false
			}
		}
		return true
	}
}

func Or(rs ...Rule) Rule {
	return func(ev atmt2.Event) bool {
		for _, r := range rs {
			if r(ev) {
				return true
			}
		}
		return false
	}
}
