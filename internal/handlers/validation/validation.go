package validation

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleBindError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errs := make(map[string]string)
		for _, e := range validationErrors {
			errs[e.Field()] = e.Tag()
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errs})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Println(err.Error())
	}
}
