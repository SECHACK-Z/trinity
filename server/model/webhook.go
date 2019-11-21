package model

import "github.com/jinzhu/gorm"

type Webhook struct {
	gorm.Model
	URL string
	body string
	Header string
	Event string
}

func (repo *GormRepository)CreateWebhook(webhook *Webhook) error {
	return repo.db.Create(webhook).Error
}

func (repo *GormRepository)GetWebhooks() (*[]Webhook, error) {
	webhooks := &[]Webhook{}
	err := repo.db.Find(webhooks).Error
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (repo *GormRepository)GetWebhooksByID(id int) (*Webhook, error) {
	webhook := &Webhook{}
	if err := repo.db.First(webhook, id).Error; err != nil {
		return nil, err
	}
	return webhook, nil
}

func (repo *GormRepository)UpdateWebhook(webhook *Webhook) error {
	return repo.db.Save(webhook).Error
}
