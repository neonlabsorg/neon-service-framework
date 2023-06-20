package api

import (
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/helpers"
	"gopkg.in/go-playground/validator.v9"
)

func HttpErrorHandler(err error, c echo.Context) {
	if err, ok := err.(*net.OpError); ok && err.Op == "write" {
		return
	}

	if ve, ok := err.(*ValidateError); ok {
		ValidationError(c, ve.Veer, ve.Model)
		return
	}

	if ve, ok := err.(errors.Error); ok {
		ProjectError(c, ve)
		return
	}

	message := err.Error()
	statusCode := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		statusCode = he.Code
		if reflect.TypeOf(he.Message).Kind().String() == "string" {
			message = he.Message.(string)
		}
	}

	if re, ok := err.(*RequestError); ok {
		statusCode = re.StatusCode()
	}

	_ = HttpError(c, statusCode, message)
}

func ProjectError(c echo.Context, err error) {

	perr := err.(errors.Error)
	message := perr.Error()
	name := perr.GetType()
	var statusCode int

	switch perr.GetType() {
	case errors.Validation, errors.Logical:
		statusCode = 400
	case errors.AccessDenied:
		statusCode = 403
	case errors.Unauthorized:
		statusCode = 401
	case errors.NotFound:
		statusCode = 404
	case errors.Temporarily:
		statusCode = 503
	default:
		statusCode = 500
	}

	_ = c.JSON(statusCode, ErrorResponseModel{
		Error: HttpErrorResponseModel{
			Message:    message,
			StatusCode: statusCode,
			Name:       name.String(),
			Context:    perr.GetContext(),
		},
	})
}

func HttpError(c echo.Context, statusCode int, message ...string) error {
	var echoError *echo.HTTPError
	if len(message) > 0 {
		echoError = echo.NewHTTPError(statusCode, message[0])
	} else {
		echoError = echo.NewHTTPError(statusCode)
	}
	return c.JSON(echoError.Code, ErrorResponseModel{
		Error: HttpErrorResponseModel{
			Message:    echoError.Message.(string),
			StatusCode: echoError.Code,
		},
	})
}

func NewRequestError(message string, statusCode ...int) error {
	re := &RequestError{message, http.StatusBadRequest}
	if len(statusCode) > 0 {
		re.statusCode = statusCode[0]
	}
	return re
}

type RequestError struct {
	message    string
	statusCode int
}

func (e *RequestError) Error() string {
	return e.message
}

func (e *RequestError) StatusCode() int {
	return e.statusCode
}

func ValidationError(c echo.Context, err error, model interface{}) {
	message := "Request is invalid"
	fields := []ValidationErrorFieldModel{}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		message = err.Error()
	} else {
		reflectType := reflect.TypeOf(model)
		if reflectType != nil && reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}

		for _, err := range err.(validator.ValidationErrors) {
			jsonName := helpers.UnderScoreString(err.Field())
			if model != nil {
				if rField, ok := reflectType.FieldByName(err.Field()); ok {
					jsonName = rField.Tag.Get("json")
				}
			}

			var message string
			if err.Param() != "" {
				message = fmt.Sprintf(`Invalid field '%s' by tag '%s = %s'`, jsonName, err.Tag(), err.Param())
			} else {
				message = fmt.Sprintf(`Invalid field '%s' by tag '%s'`, jsonName, err.Tag())
			}

			fields = append(fields, ValidationErrorFieldModel{
				FieldName: jsonName,
				Namespace: err.Namespace(),
				Tag:       err.Tag(),
				TagParam:  err.Param(),
				Message:   message,
			})
		}
	}

	_ = c.JSON(http.StatusBadRequest, ValidationErrorResponseModel{
		Error: ValidationErrorModel{
			Message:    message,
			Name:       "validation_error",
			StatusCode: http.StatusBadRequest,
			Fields:     fields,
		},
	})
}
