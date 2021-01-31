package handler

import (
	"fmt"

	"github.com/Eleron8/filestackTestApi/models"
	"github.com/labstack/echo/v4"
)

type UsecaseInter interface {
	FileFlow(data models.TransformData) error
}

type Handler struct {
	Usecase UsecaseInter
}

func NewHandler(usecase UsecaseInter) Handler {
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
		handleErr(err)
	}
	if err := h.Usecase.FileFlow(req); err != nil {
		handleErr(err)
	}
	return nil
}
