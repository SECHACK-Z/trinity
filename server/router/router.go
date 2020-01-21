package router

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"main/manager"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
)

var (
	cookieSecret = os.Getenv("COOKIE_SECRET")
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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cookieSecret))))
	e.GET("/*", r.getStaticHandler())

	e.POST("/api/login", r.postLogin)
	e.POST("/api/github", r.postGitHubWebhook)

	api := e.Group("/api")
	if cookieSecret == "" {
		cookieSecret = "secret"
	}
	api.Use(r.checkLogin)

	api.POST("/logout", r.postLogout)

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

	return nil
}

func (r *router) SetUpForInitilize(e *echo.Echo) error {
	e.GET("/*", r.getStaticHandler())
	api := e.Group("/api")
	api.POST("/init", func(c echo.Context)error {
		reqBody := &struct{
			UserName string `json:"username"`
			Password string `json:"password"`
		}{}

		c.Bind(reqBody)

		if err := r.manager.Admin.CreateAdmin(reqBody.UserName, reqBody.Password); err != nil {
			return err
		}

		go func() {
			<-c.Request().Context().Done()
			time.Sleep(500 * time.Millisecond)
			e.Shutdown(context.Background())
		}()
		return c.NoContent(http.StatusNoContent)
	})
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
