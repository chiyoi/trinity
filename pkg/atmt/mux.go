package atmt

import (
	"sort"
	"sync"
)

type Matcher interface {
	Match(ev Event) bool
	Priority() int
}

type MatcherFunc func(ev Event) bool

func (m MatcherFunc) Match(ev Event) bool { return m(ev) }
func (m MatcherFunc) Priority() int       { return 0 }

type Handler interface{ ServeEvent(ev Event) }

type HandlerFunc func(ev Event)

func (h HandlerFunc) ServeEvent(ev Event) { h(ev) }

type ServeMux struct {
	mu sync.RWMutex
	es []entry
	p  int
}

func (mux *ServeMux) Match(ev Event) bool {
	_, m := mux.Handler(ev)
	return m != nil
}
func (mux *ServeMux) Priority() int {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	return mux.p
}

func (mux *ServeMux) ServeEvent(ev Event) {
	h, _ := mux.Handler(ev)
	if h != nil {
		h.ServeEvent(ev)
	}
}

type entry struct {
	m Matcher
	h Handler
}

func (mux *ServeMux) SetPriority(p int) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.p = p
}

func appendSorted(es []entry, e entry) []entry {
	idx := sort.Search(len(es), func(i int) bool {
		return e.m.Priority() < es[i].m.Priority()
	})
	if idx == len(es) {
		return append(es, e)
	}
	es = append(es, entry{})
	copy(es[idx+1:], es[idx:])
	es[idx] = e
	return es
}
func (mux *ServeMux) Handle(matcher Matcher, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.es = appendSorted(mux.es, entry{matcher, handler})
}
func (mux *ServeMux) HandleFunc(match func(ev Event) bool, handler func(ev Event)) {
	mux.Handle(MatcherFunc(match), HandlerFunc(handler))
}

func (mux *ServeMux) Handler(ev Event) (h Handler, m Matcher) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	for _, e := range mux.es {
		if e.m.Match(ev) {
			return e.h, e.m
		}
	}
	return
}

func NewServeMux() *ServeMux { return new(ServeMux) }

var defaultServeMux ServeMux
var DefaultServeMux = &defaultServeMux
