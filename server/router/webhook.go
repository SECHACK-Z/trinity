package router

import (
	"fmt"
	"main/manager/webhook"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/webhooks.v5/github"
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

func (r *router) receiveGitHubWebook(c echo.Context) error {
	hook, _ := github.New()

	payload, err := hook.Parse(c.Request(), github.PushEvent)
	if err != nil {
		fmt.Println(err)
	}
	switch payload.(type) {
	case github.PushPayload:
		release := payload.(github.PushPayload)
		ref := strings.Split(release.Ref, "/")
		branch := ref[len(ref)-1]
		if branch == "master" {
			fmt.Println(branch)
		}
	}
	return c.String(200, "")
}
