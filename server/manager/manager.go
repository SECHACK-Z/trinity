package manager

import (
	"github.com/jinzhu/gorm"
	"main/manager/config"
	"main/manager/healthcheck"
	"main/manager/webhook"
)

type Manager struct {
	Config  *config.ConfigManager
	Webhook *webhook.WebhookManager
	HealthCheck *healthcheck.HealthCheckManager
}

func New(db *gorm.DB) *Manager {
	return &Manager{
		Config:  config.New(db),
		Webhook: webhook.New(db),
		HealthCheck: healthcheck.New(db),
	}
}
