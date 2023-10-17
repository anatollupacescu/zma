package main

import (
	"github.com/anatollupacescu/zma/bmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &app{
		collections: make(map[string]*bmt.Bmt),
	}

	e.POST("/:collection", app.upload)
	e.GET("/proof/:collection/:index", app.proof)

	e.Logger.Fatal(e.Start(":8080"))
}
