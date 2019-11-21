package router

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
	"github.com/ymotongpoo/goltsv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"main/config"
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

	e.GET("/api/log", func(c echo.Context) error {
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
	})
	e.GET("/api/alp", func(c echo.Context) error {
		out, err := exec.Command("sh", "-c", "cat logFile | alp --sort=max ltsv").Output()
		if err != nil {
			pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: err.Error()})
		}
		return c.String(200, string(out))
	})

	e.GET("/api/config", func(c echo.Context) error {
		config := r.configManager.Get()
		buf, err := yaml.Marshal(config)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, struct {
			Yaml string
		}{
			Yaml: string(buf),
		})
	})

	e.POST("/api/config", func(c echo.Context) error {
		req := &struct {
			Yaml string `json:"yaml"`
		}{}
		c.Bind(req)

		newConfig := config.Config{}
		if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if err := r.configManager.SetConfig(newConfig); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusOK)
	})

	e.POST("/api/config/save", func(c echo.Context) error {
		req := &struct {
			Name string `json:"name"`
			Yaml string `json:"yaml"`
		}{}
		c.Bind(req)

		newConfig := config.Config{}
		if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if err := r.configManager.SetConfig(newConfig); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := r.configManager.Save(); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
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
