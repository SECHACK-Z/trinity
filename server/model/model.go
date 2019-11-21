package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	CreateWebhook(*Webhook) error
	GetWebhooks() (*[]Webhook, error)
	GetWebhooksByID(int) (*Webhook, error)
	UpdateWebhook(*Webhook) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(engine *gorm.DB) Repository {
	return &GormRepository{db: engine}
}

