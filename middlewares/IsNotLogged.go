package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
)

func IsNotLogged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.(*CustomContext).Auth()
		if auth.Authenticated == false {
			return next(c)
		}
		return c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}
