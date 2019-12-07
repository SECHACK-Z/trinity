package router

import (
	"fmt"
	"main/manager/webhook"
	"main/pubsub"
	"main/pubsub/systemevent"
	"time"

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
	pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.WEBHOOK_RECEIVED})

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
		URL := release.Repository.URL

		for _, target := range r.manager.Config.Get().Targets {
			if target.Repository == URL && "master" == branch {
				message := fmt.Sprintf("New commit is pushed on %s at %s\n", branch, URL)
				pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.REPOSITORY_UPDATED, Message: message})
				pubsub.GetWebookEvent.Pub(pubsub.GetWebook{Repository: URL})
				return c.NoContent(200)
			}
		}

	}
	pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.REPOSITORY_UPDATED, Message: "Webhook was came but no settings found."})
	return c.NoContent(404)
}
