package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Eleron8/filestackTestApi/configuration"
	"github.com/Eleron8/filestackTestApi/getfile"
	"github.com/Eleron8/filestackTestApi/gstorage"
	"github.com/Eleron8/filestackTestApi/handler"
	"github.com/Eleron8/filestackTestApi/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	fHandl := getfile.NewFileHandler(httpClient)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	gStorage, err := gstorage.NewStorageData(ctx, configuration.Config.ProjectID, configuration.Config.BucketName)
	if err != nil {
		configuration.Logger.Fatal("can't connect to GCS", zap.Error(err))
	}
	useCase := usecase.NewUsecase(fHandl, gStorage, configuration.Config.FolderName, configuration.Config.MaxGoroutines, configuration.Logger)
	routeHandler := handler.NewHandler(useCase)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("accept", routeHandler.Accept)
	e.Logger.Fatal(e.Start(":" + configuration.Config.ServerPort))
}
