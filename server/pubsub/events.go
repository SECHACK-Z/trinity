package pubsub

import (
	"main/config"
	"main/pubsub/systemevent"

	"net/http"
	"time"
)

type Access struct {
	Req     *http.Request
	Res     *http.Response
	Elapsed time.Duration
}

type System struct {
	Time    time.Time
	Type    systemevent.SystemEventType
	Message string
}

type UpdateConfig struct {
	Config config.Config
}

type HealthCheck struct {
	Target  config.Target
	Status  int
	Message string
}

type GetWebhook struct {
	Repository string
}

type News struct {
	Type   string
	Target config.Target
}
