package router

import (
	"main/manager"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
)

type router struct {
	manager *manager.Manager
}

func New(manager *manager.Manager) *router {
	return &router{
		manager: manager,
	}
}

func (r *router) SetUp(e *echo.Echo) error {
	e.GET("/*", r.getStaticHandler())

	api := e.Group("/api")

	api.GET("/404", r.defaultBackend)
	api.GET("/ping", r.ping)
	api.GET("/logs", r.getRawLogs)

	api.GET("/config", r.getConfig)
	api.POST("/config", r.postConfig)
	api.POST("/config/save", r.postSaveConfig)

	api.GET("/webhooks", r.getWebhooks)
	api.POST("/webhooks", r.postWebhooks)
	api.PUT("/webhooks/:id", r.putWebhookByID)
	api.DELETE("/webhooks/:id", r.deleteWebhookByID)

	api.POST("/github", r.receiveGitHubWebook)
	return nil
}

func (r *router) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (r *router) getStaticHandler() func(c echo.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		panic(err)
	}

	h := http.FileServer(statikFs)
	return echo.WrapHandler(http.StripPrefix("/", h))
}

func (r *router) defaultBackend(c echo.Context) error {
	return c.HTML(200, "<h1>Welcome to <かっこいい名前></h1>")
}
