package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func LogoutMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.(*CustomContext).Auth()
		if auth.Authenticated == true {
			session, _ := session.Get("session", c)
			session.Values["authenticated"] = false
			session.Values["name"] = ""
			session.Values["lastname"] = ""
			session.Values["email"] = ""
			session.Save(c.Request(), c.Response())
			return next(c)
		}
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}
