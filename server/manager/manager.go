package manager

import (
	"github.com/jinzhu/gorm"
	"main/manager/config"
	"main/manager/webhook"
)

type Manager struct {
	Config *config.ConfigManager
	Webhook *webhook.WebhookManager
}

func New(db *gorm.DB)*Manager {
	return &Manager{
		Config:  config.New(db),
		Webhook: webhook.New(db),
	}
}
