package webhook

import "github.com/jinzhu/gorm"

type Webhook struct {
	gorm.Model
	URL     string          `json:"url"`
	Body    string          `json:"body"`
	Header  string          `json:"header"`
	Event   []WebhookEvent  `json:"event"`
	Secrets []WebhookSecret `json:"secrets"`
}

type WebhookEvent struct {
	gorm.Model
	WebhookID uint   `json:"webhook_id"`
	Event     string `json:"event"`
}

type WebhookSecret struct {
	gorm.Model
	WebhookID   uint   `json:"webhook_id"`
	PlaceHolder string `json:"place_holder"`
	Secret      string `json:"secret"`
}
