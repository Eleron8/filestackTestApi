package main

import (
	"net/http"
	"time"

	"github.com/Eleron8/filestackTestApi/getfile"
	"github.com/Eleron8/filestackTestApi/handler"
	"github.com/Eleron8/filestackTestApi/usecase"
	"github.com/labstack/echo/v4"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	fHandl := getfile.NewFileHandler(httpClient)
	useCase := usecase.NewUsecase(fHandl)
	routeHandler := handler.NewHandler(useCase)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("accept", routeHandler.Accept)
	e.Logger.Fatal(e.Start(":8080"))
}
