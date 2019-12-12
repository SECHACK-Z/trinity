package manager

import (
	"github.com/jinzhu/gorm"
	"main/manager/cd"
	"main/manager/config"
	"main/manager/healthcheck"
	"main/manager/log"
	"main/manager/webhook"
)

type Manager struct {
	Config             *config.ConfigManager
	Webhook            *webhook.WebhookManager
	HealthCheck        *healthcheck.HealthCheckManager
	ContinuousDelivery *cd.CDManager
	Log                *log.LogManager
}

func New(db *gorm.DB) *Manager {
	return &Manager{
		Config:             config.New(db),
		Webhook:            webhook.New(db),
		HealthCheck:        healthcheck.New(db),
		ContinuousDelivery: cd.New(db),
		Log:                log.New(db),
	}
}
