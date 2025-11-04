package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation")
	ErrUnauthorized = errors.New("unauthorized")
)

func Translate(err error) *echo.HTTPError {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, ErrValidation):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrUnauthorized):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrConflict):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
}
