package middlewares

import (
	"log"
	"net/http"

	. "github.com/alegrecode/echo/LoginMongoDB/models"
	"github.com/gookit/validate"
	"github.com/labstack/echo"
)

func ValidateRegister(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := new(User)
		if err2 := c.Bind(data); err2 != nil {
			log.Fatal(err2)
		}

		v := validate.Struct(data)
		if !v.AtScene("register").Validate() {
			c.(*CustomContext).SetFlash("error", v.Errors)
			return c.Redirect(http.StatusMovedPermanently, "/register")
		}
		m := map[string]interface{}{
			"password":         c.FormValue("password"),
			"confirm_password": c.FormValue("confirm_password"),
		}
		vm := validate.Map(m)
		vm.AddRule("confirm_password", "eqField", "password").SetMessage("Fields password and confirm password not match.")
		if !vm.Validate() {
			c.(*CustomContext).SetFlash("error", vm.Errors)
			return c.Redirect(http.StatusMovedPermanently, "/register")
		}
		return next(c)
	}
}
