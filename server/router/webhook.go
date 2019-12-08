package router

import (
	"main/manager/webhook"

	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *router) getWebhooks(c echo.Context) error {
	webhooks, err := r.manager.Webhook.GetWebhooks()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, webhooks)
}

func (r *router) postWebhooks(c echo.Context) error {
	webhook := &webhook.Webhook{}
	c.Bind(webhook)
	if err := r.manager.Webhook.CreateWebhook(webhook); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, webhook)
}

func (r *router) putWebhookByID(c echo.Context) error {
	webhook := &webhook.Webhook{}
	c.Bind(webhook)
	if err := r.manager.Webhook.UpdateWebhook(webhook); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (r *router) deleteWebhookByID(c echo.Context) error {
	webhook := &webhook.Webhook{}
	c.Bind(webhook)
	if err := r.manager.Webhook.DeleteWebhook(webhook); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)

}
