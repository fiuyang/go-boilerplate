package helper

import (
	"fmt"
	"gin-boilerplate/data/response"
	"gin-boilerplate/model"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var modelMap = map[string]reflect.Type{
	"users": reflect.TypeOf(model.Users{}),
	"tags":  reflect.TypeOf(model.Tags{}),
}

func ValidateUnique(db *gorm.DB, fl validator.FieldLevel) bool {
	value := fl.Field().String()
	tableName := getModelFromTag(fl)
	fmt.Println("Validation: EmailExistsInTable for table", tableName)

	// Pass the actual db instance to EmailExistsInTable
	exists := EmailExistsInTable(db, value, tableName)

	return !exists
}

func EmailExistsInTable(db *gorm.DB, value, tableName string) bool {
	modelType, ok := modelMap[tableName]
	if !ok {
		fmt.Printf("Unknown model name: %s\n", tableName)
		return false
	}

	// Create a new instance of the dynamic model
	modelInstance := reflect.New(modelType).Interface()
    fmt.Println(value)
	if err := db.Table(tableName).Where("email = ?", value).First(modelInstance).Error; err != nil {
		return false
	}
	return true
}


func getModelFromTag(fl validator.FieldLevel) string {
	// Assuming 'validate' tag is in the format "unique=tableName"
	validateTag := fl.Param()

	parts := strings.Split(validateTag, "=")
    if len(parts) >= 2 {
        return parts[1]
    }

    return parts[0]

}

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
				errorMap[e.Field()] = fmt.Sprintf("%s tidak valid", fieldName)
			case "gte":
				errorMap[e.Field()] = fmt.Sprintf("%s harus lebih besar dari %s", fieldName, e.Param())
			case "lte":
				errorMap[e.Field()] = fmt.Sprintf("%s harus lebih kecil dari %s", fieldName, e.Param())
			case "unique":
				errorMap[e.Field()] = fmt.Sprintf("%s has already been taken in table %s", fieldName, e.Param())
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
