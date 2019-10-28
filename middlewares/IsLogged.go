package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
)

func IsLogged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.(*CustomContext).Auth()
		if auth.Authenticated == true {
			return next(c)
		}
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}
