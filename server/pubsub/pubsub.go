package pubsub

import (
	"github.com/cheekybits/genny/generic"
)

type EventType generic.Type

type __EventTypePubSub struct {
	subs map[string]func(EventType)
	c    chan EventType
}

var EventTypeEvent = &__EventTypePubSub{
	subs: make(map[string]func(EventType)),
	c:    make(chan EventType, 10),
}

func (ps *__EventTypePubSub) Sub(f func(et EventType)) string {
	subID := randomStr(5)
	for _, ok := ps.subs[subID]; ok; _, ok = ps.subs[subID] {
		subID = randomStr(5)
	}
	ps.subs[subID] = f
	return subID
}

func (ps *__EventTypePubSub) Unsub(subscribeID string) bool {
	if _, ok := ps.subs[subscribeID]; ok {
		delete(ps.subs, subscribeID)
		return true
	}
	return false
}

func (ps *__EventTypePubSub) Pub(event EventType) {
	for _, f := range ps.subs {
		go f(event)
	}
}


