package router

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
	"github.com/ymotongpoo/goltsv"
	"io/ioutil"
	"main/config/manager"
	"main/model"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http"
	"os/exec"
	"time"
)

type router struct {
	repo          model.Repository
	configManager *manager.ConfigManager
}

func New(repo model.Repository, configManager *manager.ConfigManager) *router {
	return &router{
		repo:          repo,
		configManager: configManager,
	}
}

func (r *router) SetUp(e *echo.Echo) error {
	e.GET("/api/ping", r.ping)
	e.GET("/*", r.getStaticHandler())

	e.GET("/api/log", r.getLog)
	e.GET("/api/alp", r.getALP)
	e.GET("/api/config", r.getConfig)
	e.POST("/api/config", r.postConfig)
	e.POST("/api/config/save", r.postSaveConfig)
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

func (r *router) getLog(c echo.Context) error {
	//デモ用の実装
	data, err := ioutil.ReadFile(`logFile`)
	if err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: err.Error()})
	}
	b := bytes.NewBufferString(string(data))
	reader := goltsv.NewReader(b)
	records, _ := reader.ReadAll()
	bytes, _ := json.Marshal(records)
	return c.JSON(200, string(bytes))
}

func(r *router) getALP (c echo.Context) error {
	out, err := exec.Command("sh", "-c", "cat logFile | alp --sort=max ltsv").Output()
	if err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: err.Error()})
	}
	return c.String(200, string(out))
}

