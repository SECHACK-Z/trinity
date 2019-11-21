package pubsub

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "EventType=`cat events.go | grep struct | awk '{print \$2}' | paste -s  -d , -`"

import (
	"github.com/cheekybits/genny/generic"
)

type EventType generic.Type

type __EventTypePubSub struct {
	subs map[string]func(EventType)
	c chan EventType
}

var _EventTypePubSub *__EventTypePubSub

func GetEventTypePubSub() *__EventTypePubSub {
	if _EventTypePubSub == nil {
		_EventTypePubSub = &__EventTypePubSub{
			subs: make(map[string]func(EventType)),
			c: make(chan EventType, 10),
		}
	}
	return _EventTypePubSub
}

func (ps *__EventTypePubSub)Sub(f func(et EventType))string {
	subID := randomStr(5)
	for _, ok := ps.subs[subID];ok; _, ok = ps.subs[subID] {
		subID = randomStr(5)
	}
	ps.subs[subID] = f
	return subID
}

func (ps *__EventTypePubSub)Pub(event EventType) {
	for _, f := range ps.subs {
		go f(event)
	}
}
