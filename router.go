package main

import (
	"log"
	"mfbo_api/config"
	"mfbo_api/controllers"
	"mfbo_api/db"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func router() *echo.Echo {
	db.Db = db.DBInit()
	db.DbDashboard = db.DBDashboardInit()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.CORSAllowOrigin,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Use(printUrlMiddleware)
	e.Use(middleware.Logger())

	//admin := e.Group("/admin")

	//home
	e.GET("/", controllers.HelloGreeting)

	//admin.Use(lib.AuthenticationMiddleware)

	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}
