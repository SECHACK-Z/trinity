package pubsub

import (
	"github.com/cheekybits/genny/generic"
)

type EventType generic.Type

// TODO: Close処理
type EventTypePubSub struct {
	subs map[string]func(EventType)
	c chan EventType
}

var _EventTypePubSub *EventTypePubSub

func GetEventTypePubSub() *EventTypePubSub {
	if _EventTypePubSub == nil {
		_EventTypePubSub = &EventTypePubSub{
			subs: make(map[string]func(EventType)),
			c: make(chan EventType, 10),
		}
	}
	return _EventTypePubSub
}

func (ps *EventTypePubSub)Sub(f func(et EventType))string {
	subID := randomStr(5)
	for _, ok := ps.subs[subID];ok; _, ok = ps.subs[subID] {
		subID = randomStr(5)
	}
	ps.subs[subID] = f
	return subID
}

func (ps *EventTypePubSub)Pub(event EventType) {
	for _, f := range ps.subs {
		go f(event)
	}
}
