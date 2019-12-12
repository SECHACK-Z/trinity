package router

import (
	"bytes"
	"io/ioutil"
	"main/pubsub"
	"main/pubsub/systemevent"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ymotongpoo/goltsv"
)

func (r *router) getRawLogs(c echo.Context) error {
	data, err := ioutil.ReadFile(`./accessLog`)
	if err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: err.Error()})
	}
	b := bytes.NewBufferString(string(data))

	reader := goltsv.NewReader(b)
	records, _ := reader.ReadAll()
	return c.JSON(200, records) //内部的にJSON.marshalしているらしい
}
