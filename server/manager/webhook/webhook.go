package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http"
	"strings"
	"time"
)

type WebhookManager struct {
	db *gorm.DB
	statusMap map[string]int
}

func New(db *gorm.DB) *WebhookManager {
	db.AutoMigrate(&Webhook{}, &WebhookEvent{})
	webhookManager := &WebhookManager{
		db: db,
		statusMap: make(map[string]int),
	}
	pubsub.UpdateConfigEvent.Sub(webhookManager.onUpdateConfig)
	pubsub.HealthCheckEvent.Sub(webhookManager.onHealthCheck)
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
	message := "新しいコンフィグが適用されました"
	for _, webhook := range webhooks {
		go callWebhook(webhook, message)
	}
}

func (m *WebhookManager) onHealthCheck(event pubsub.HealthCheck) {
	pre, ok := m.statusMap[event.Target]
	if !ok {
		pre = 200
	}
	if pre < 400 && 400 <= event.Status {
		webhooks, err := m.GetWebhooksByEvent("po")
		if err != nil {
			pubsub.SystemEvent.Pub(pubsub.System{
				Time:    time.Now(),
				Type:    systemevent.ERROR,
				Message: err.Error(),
			})
		}
		message := event.Target + " のヘルスチェックに失敗しました"
		for _, webhook := range webhooks {
			go callWebhook(webhook, message)
		}
	}

	if pre > 400 && 400 > event.Status {
		webhooks, err := m.GetWebhooksByEvent("po")
		if err != nil {
			pubsub.SystemEvent.Pub(pubsub.System{
				Time:    time.Now(),
				Type:    systemevent.ERROR,
				Message: err.Error(),
			})
		}
		message := event.Target + " が回復しました"
		for _, webhook := range webhooks {
			go callWebhook(webhook, message)
		}
	}

	m.statusMap[event.Target] = event.Status
}

func callWebhook(webhook *Webhook, message string) (*http.Response, error) {
	body := strings.ReplaceAll(webhook.Body, "<message>", message)
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewReader([]byte(body)))
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
	m.db.Model(&Webhook{}).Association("Event").Replace(webhook.Event)
	return m.db.Save(webhook).Error
}

func (m *WebhookManager) DeleteWebhook(webhook *Webhook) error {
	return m.db.Delete(webhook).Error
}
