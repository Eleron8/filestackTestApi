package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Eleron8/filestackTestApi/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UsecaseInter interface {
	FileFlow(data models.TransformData, wr io.Writer) error
}

type Handler struct {
	Usecase UsecaseInter
	logger  *zap.Logger
}

func NewHandler(usecase UsecaseInter, logger *zap.Logger) Handler {
	return Handler{
		Usecase: usecase,
	}
}

func (h Handler) Accept(ctx echo.Context) error {
	handleErr := func(err error) error {
		return fmt.Errorf("accept request to transform: %w", err)
	}
	ctx.Response().Header().Set(echo.HeaderContentType, "application/zip")
	ctx.Response().WriteHeader(http.StatusOK)

	var req models.TransformData
	if err := ctx.Bind(&req); err != nil {
		h.logger.Info("can't bind body", zap.Error(err))
		return handleErr(err)
	}
	f, err := os.Create("archive.zip")
	if err != nil {
		return handleErr(err)
	}
	// newFile, err := os.Open("archive.zip")
	// if err != nil {
	// 	return handleErr(err)
	// }

	wr := ctx.Response().Writer

	err = h.Usecase.FileFlow(req, f)
	if err != nil {
		h.logger.Info("image transformation flow failed", zap.Error(err))
		return handleErr(err)
	}
	newfile, err := ioutil.ReadFile("archive.zip")
	if err != nil {
		return handleErr(err)
	}
	wr.Write(newfile)
	ctx.Response().Flush()

	// return ctx.Stream(http.StatusOK, "application/zip", newFile)
	return nil
}
