package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

func (r *router) postLogin(c echo.Context) error {
	reqBody := &struct{
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(reqBody)

	adminUser, err := r.manager.Admin.Authorize(reqBody.UserName, reqBody.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	sess.Options = &sessions.Options{
		Path: "/",
		MaxAge: 86400 * 7,
		HttpOnly: true,
	}
	sess.Values["username"] = adminUser.UserName
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusNoContent)
}

func (r *router) postLogout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Values["username"] = nil
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusNoContent)
}

func (r *router)checkLogin(next echo.HandlerFunc)echo.HandlerFunc {
	return func (c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		if sess.Values["username"] == nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}
