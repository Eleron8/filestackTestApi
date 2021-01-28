package handler

import (
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Accept(ctx echo.Context) error {

	return nil
}
