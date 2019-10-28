package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	. "github.com/alegrecode/echo/LoginMongoDB/middlewares"
	"github.com/alegrecode/echo/LoginMongoDB/models"
)

func LoginView(c echo.Context) error {
	getFlashes := c.(*CustomContext).GetFlash()
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"flash": getFlashes,
	})
}

func RegisterView(c echo.Context) error {
	getFlashes := c.(*CustomContext).GetFlash()
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"flash": getFlashes,
	})
}

func RegisterUser(c echo.Context) error {
	result := models.SaveUser(c)
	c.(*CustomContext).SetFlash("done", "User saved success. Now you can sign in.")
	fmt.Println(result.InsertedID)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func AuthUser(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/dashboard")
}

func Dashboard(c echo.Context) error {
	cc := c.(*CustomContext)
	auth := cc.Auth()
	return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
		"auth": auth,
	})
}

func Logout(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/")
}
