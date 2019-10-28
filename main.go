package main

import (
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/alegrecode/echo/LoginMongoDB/controllers"
	"github.com/alegrecode/echo/LoginMongoDB/db"
	. "github.com/alegrecode/echo/LoginMongoDB/helpers"
	. "github.com/alegrecode/echo/LoginMongoDB/middlewares"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("assets"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("alegrecode"))))

	templates := make(map[string]*template.Template)
	templates["login.html"] = template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseFiles("views/login.html", "views/base.html", "views/navbar.partial.html", "views/alert.partial.html"))
	templates["register.html"] = template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseFiles("views/register.html", "views/base.html", "views/navbar.partial.html", "views/alert.partial.html"))
	templates["dashboard.html"] = template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseFiles("views/dashboard.html", "views/base.html", "views/navbar.partial.html", "views/alert.partial.html"))

	e.Renderer = &TemplateRegistry{
		Templates: templates,
	}

	c := db.GetClient()
	err = c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}
	DB := db.GetDatabase()
	fmt.Println("MongoDB connected " + DB.Name())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return next(cc)
		}
	})

	e.GET("/", controllers.LoginView, IsNotLogged)

	e.GET("/register", controllers.RegisterView, IsNotLogged)

	e.POST("/register", controllers.RegisterUser, ValidateRegister)

	e.POST("/login", controllers.AuthUser, ValidateLogin)

	e.GET("/dashboard", controllers.Dashboard, IsLogged)

	e.DELETE("/logout", controllers.Logout, LogoutMiddleware)

	e.Logger.Fatal(e.Start(":8080"))

}
