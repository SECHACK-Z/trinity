package pubsub

import (
	"main/pubsub/systemevent"
	"net/http"
	"time"
)

type AccessEvent struct {
	Req *http.Request
	Res *http.Response
	Elapsed time.Duration
}

type SystemEvent struct {
	Time time.Time
	Type systemevent.SystemEventType
	Message string
}
