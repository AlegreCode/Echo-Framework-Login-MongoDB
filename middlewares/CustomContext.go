package middlewares

import (
	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type CustomContext struct {
	echo.Context
}

type Auth struct {
	Name          string `json:"name"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	Authenticated bool   `json:"authenticated"`
}

func (c CustomContext) Auth() Auth {
	var auth Auth
	session, _ := session.Get("session", c)
	if session.Values["authenticated"] == nil {
		auth.Name = ""
		auth.Lastname = ""
		auth.Email = ""
		auth.Authenticated = false
		return auth
	}
	auth.Name = session.Values["name"].(string)
	auth.Lastname = session.Values["lastname"].(string)
	auth.Email = session.Values["email"].(string)
	auth.Authenticated = session.Values["authenticated"].(bool)
	return auth
}

func (c CustomContext) SetFlash(t string, msg interface{}) bool {
	session, _ := session.Get("flash", c)
	session.Options = &sessions.Options{
		MaxAge: 1,
	}
	flashes := map[string]interface{}{"type": t, "msg": msg}
	flashesJSON, _ := json.Marshal(flashes)
	flashesString := string(flashesJSON)
	session.AddFlash(flashesString)
	session.Save(c.Request(), c.Response())
	return true
}

func (c CustomContext) GetFlash() map[string]interface{} {
	session, _ := session.Get("flash", c)
	if flash := session.Flashes(); len(flash) > 0 {
		flashes := make(map[string]interface{})
		json.Unmarshal([]byte(flash[0].(string)), &flashes)
		return flashes
	}
	return map[string]interface{}{}
}
