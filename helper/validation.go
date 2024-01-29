package helper

import (
	"fmt"
	"gin-boilerplate/data/response"
	"net/http"
	"reflect"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationError(err error, c *gin.Context, requestType interface{}) bool {
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		errorMap := make(map[string]string)
		
		for _, e := range castedObject {
			field, _ := reflect.TypeOf(requestType).FieldByName(e.StructField())
			fieldName := field.Tag.Get("json")
			switch e.Tag() {
			case "required":
				errorMap[e.Field()] = fmt.Sprintf("%s wajib diisi", fieldName)
			case "email":
				errorMap[e.Field()] = fmt.Sprintf("%s bukan email yang valid", fieldName)
			case "gte":
				errorMap[e.Field()] = fmt.Sprintf("%s harus lebih besar dari %s", fieldName, e.Param())
			case "lte":
				errorMap[e.Field()] = fmt.Sprintf("%s harus lebih kecil dari %s", fieldName, e.Param())
			}
		}
		
		c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Status:  "BadRequest",
			Message: "Error Validation",
			Errors: errorMap,
		})
		return true
	}
	return false
}
