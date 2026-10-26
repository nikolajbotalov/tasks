package validation

import (
	"TaskFlow/internal/delivery/response"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" || name == "-" {
				return fld.Name
			}
			return name
		})
	}
}

func HandleBindError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make([]response.FieldError, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			details = append(details, response.FieldError{
				Field:   fieldErr.Field(),
				Message: validationMessage(fieldErr),
			})
		}

		c.JSON(http.StatusUnprocessableEntity, response.ValidationErrorResponse{
			Error:   "validation failed",
			Details: details,
		})
		return
	}

	c.JSON(http.StatusBadRequest, response.ErrorResponse{
		Error: response.InvalidRequestBody,
	})
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email"
	case "uuid":
		return "must be a valid UUID"
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters", err.Param())
		}
		return fmt.Sprintf("must be at least %s", err.Param())
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters", err.Param())
		}
		return fmt.Sprintf("must be at most %s", err.Param())
	case "oneof":
		return "must be one of: " + strings.Join(strings.Fields(err.Param()), ", ")
	default:
		return "invalid value"
	}
}
