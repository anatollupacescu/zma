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
		collections: make(map[string]*bmt.Bmt[[]byte]),
	}

	e.POST("/:collection", app.upload)
	e.GET("/:collection/:index", app.download)
	e.GET("/proof/:collection/:index", app.proof)

	e.Logger.Fatal(e.Start(":8080"))
}
