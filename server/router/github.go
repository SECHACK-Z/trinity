package router

import (
	"fmt"
	"main/pubsub"
	"main/pubsub/systemevent"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func (r *router) postGitHubWebhook(c echo.Context) error {
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
				pubsub.DeployEvent.Pub(pubsub.Deploy{Repository: URL, Type: "がデプロイされました"})
				pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.REPOSITORY_UPDATED, Message: message})
				pubsub.GetWebhookEvent.Pub(pubsub.GetWebhook{Repository: URL})
				return c.NoContent(200)
			}
		}

	}
	pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.REPOSITORY_UPDATED, Message: "Webhook was came but no settings found."})
	return c.NoContent(404)
}
