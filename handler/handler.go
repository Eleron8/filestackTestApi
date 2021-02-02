package handler

import (
	"fmt"
	"net/http"

	"github.com/Eleron8/filestackTestApi/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UsecaseInter interface {
	FileFlow(data models.TransformData) error
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
	var req models.TransformData
	if err := ctx.Bind(&req); err != nil {
		h.logger.Info("can't bind body", zap.Error(err))
		handleErr(err)
	}
	if err := h.Usecase.FileFlow(req); err != nil {
		h.logger.Info("image transformation flow failed", zap.Error(err))
		handleErr(err)
	}
	return ctx.NoContent(http.StatusOK)
}
