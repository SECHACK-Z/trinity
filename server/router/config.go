package router

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
	"main/config"
	"net/http"
)

func (r *router) getConfig(c echo.Context) error {
	config := r.manager.Config.Get()
	buf, err := yaml.Marshal(config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, struct {
		Yaml string
	}{
		Yaml: string(buf),
	})
}

func (r *router) postConfig(c echo.Context) error {
	req := &struct {
		Yaml string `json:"yaml"`
	}{}
	c.Bind(req)

	newConfig := config.Config{}
	if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := r.manager.Config.SetConfig(newConfig); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

func (r *router) postSaveConfig(c echo.Context) error {
	req := &struct {
		Name string `json:"name"`
		Yaml string `json:"yaml"`
	}{}
	c.Bind(req)

	newConfig := config.Config{}
	if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := r.manager.Config.SetConfig(newConfig); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := r.manager.Config.Save(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
