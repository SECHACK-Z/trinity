package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http"
	"time"
	"fmt"
)


type WebhookManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) *WebhookManager {
	db.AutoMigrate(&Webhook{}, &WebhookEvent{})
	webhookManager := &WebhookManager{db: db}
	pubsub.UpdateConfigEvent.Sub(webhookManager.onUpdateConfig)
	return webhookManager
}

type Webhook struct {
	gorm.Model
	URL    string         `json:"url"`
	Body   string         `json:"body"`
	Header string         `json:"header"`
	Event  []WebhookEvent `json:"event"`
}

type WebhookEvent struct {
	gorm.Model
	WebhookID uint   `json:"webhook_id"`
	Event     string `json:"event"`
}

func (m *WebhookManager) onUpdateConfig(event pubsub.UpdateConfig) {
	webhooks, err := m.GetWebhooksByEvent("poi")

	if err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{
			Time:    time.Now(),
			Type:    systemevent.ERROR,
			Message: err.Error(),
		})
	}
	for _, webhook := range webhooks {
		go callWebhook(webhook)
	}
}

func callWebhook(webhook *Webhook) (*http.Response, error) {
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewReader([]byte(webhook.Body)))
	if err != nil {
		return nil, err
	}
	var headers map[string]interface{}
	json.Unmarshal([]byte(webhook.Header), &headers)
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	res, err := http.DefaultClient.Do(req)
	return res, err
}

func (m *WebhookManager) CreateWebhook(webhook *Webhook) error {
	return m.db.Create(webhook).Error
}

func (m *WebhookManager) GetWebhooks() ([]*Webhook, error) {
	webhooks := []*Webhook{}
	err := m.db.Find(&webhooks).Error
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		err := m.db.Model(webhook).Related(&webhook.Event).Error
		if err != nil {
			return nil, err
		}
	}

	return webhooks, nil
}

func (m *WebhookManager) GetWebhooksByID(id int) (*Webhook, error) {
	webhook := &Webhook{}
	if err := m.db.First(webhook, id).Error; err != nil {
		return nil, err
	}

	if err := m.db.Model(webhook).Related(&webhook.Event).Error; err != nil {
		return nil, err
	}

	return webhook, nil
}

func (m *WebhookManager) GetWebhooksByEvent(eventName string) ([]*Webhook, error) {
	var webhooks []*Webhook
	if err := m.db.Joins("LEFT JOIN webhook_events ON webhooks.id = webhook_events.webhook_id").Where("webhook_events.event = ?", eventName).Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (m *WebhookManager) UpdateWebhook(webhook *Webhook) error {
	return m.db.Save(webhook).Error
}

func (m *WebhookManager) DeleteWebhook(webhook *Webhook) error {
	return m.db.Delete(webhook).Error
}
