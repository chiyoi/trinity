package atmt

import (
	"io"
	"sort"
	"sync"
)

type Matcher struct {
	Match     func(msg Message) bool
	Priority  int
	Temporary bool
}

type Handler interface {
	ServeMessage(resp *Message, post Message)
}

type HandlerFunc func(*Message, Message)

func (h HandlerFunc) ServeMessage(resp *Message, post Message) { h(resp, post) }

type ServeMux struct {
	mu sync.RWMutex
	es []entry
}

type entry struct {
	m Matcher
	h Handler
}

var _ Handler = (*ServeMux)(nil)

func (mux *ServeMux) Matcher() Matcher {
	return Matcher{
		Match: func(msg Message) bool {
			return mux.handler(msg) != -1
		},
	}
}

func (mux *ServeMux) handler(msg Message) int {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	for i, e := range mux.es {
		if e.m.Match(msg) {
			return i
		}
	}
	return -1
}

func (mux *ServeMux) ServeMessage(resp *Message, post Message) {
	idx := mux.handler(post)
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if idx != -1 {
		e := mux.es[idx]
		e.h.ServeMessage(resp, post)
		if e.m.Temporary {
			mux.es = append(mux.es[:idx], mux.es[idx+1:]...)
		}
	}
}

func appendSorted(es []entry, e entry) []entry {
	idx := sort.Search(len(es), func(i int) bool {
		return e.m.Priority < es[i].m.Priority
	})
	if idx == len(es) {
		return append(es, e)
	}
	es = append(es, entry{})
	copy(es[idx+1:], es[idx:])
	es[idx] = e
	return es
}

func (mux *ServeMux) Handle(m Matcher, h Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.es = appendSorted(mux.es, entry{m, h})
}

func (mux *ServeMux) Handler(msg Message) (m Matcher, h Handler) {
	idx := mux.handler(msg)
	if idx == -1 {
		return
	}
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	return mux.es[idx].m, mux.es[idx].h
}

func NewServeMux() *ServeMux { return new(ServeMux) }

var DefaultServeMux = NewServeMux()

type multiHandler struct {
	handlers []Handler
}

var _ Handler = (*multiHandler)(nil)

func (mh *multiHandler) ServeMessage(resp *Message, post Message) {
	for _, h := range mh.handlers {
		io.MultiWriter()
		h.ServeMessage(resp, post)
	}
}

func MultiHandler(handlers []Handler) Handler {
	return &multiHandler{handlers}
}
