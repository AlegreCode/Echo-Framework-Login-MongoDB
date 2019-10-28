package middlewares

import (
	"log"
	"net/http"

	. "github.com/alegrecode/echo/LoginMongoDB/models"
	"github.com/gookit/validate"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
)

func ValidateLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := session.Get("session", c)

		data := new(User)
		if err2 := c.Bind(data); err2 != nil {
			log.Fatal(err2)
		}
		v := validate.Struct(data)
		if !v.AtScene("login").Validate() {
			c.(*CustomContext).SetFlash("error", v.Errors)
			return c.Redirect(http.StatusMovedPermanently, "/")
		}

		email := c.FormValue("email")
		password := c.FormValue("password")
		user, err := GetSingleUser(email)
		if err != nil {
			c.(*CustomContext).SetFlash("error", err.Error())
			return c.Redirect(http.StatusMovedPermanently, "/")
		}

		if err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err1 != nil {
			c.(*CustomContext).SetFlash("error", "Email and/or password are incorrects.")
			return c.Redirect(http.StatusMovedPermanently, "/")
		}
		session.Values["authenticated"] = true
		session.Values["name"] = user.Name
		session.Values["lastname"] = user.Lastname
		session.Values["email"] = user.Email
		session.Save(c.Request(), c.Response())
		return next(c)
	}
}
