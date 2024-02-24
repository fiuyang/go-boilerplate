package helper

import (
	"fmt"
	"gin-boilerplate/model"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var modelMap = map[string]reflect.Type{
	"users": reflect.TypeOf(model.Users{}),
	"tags":  reflect.TypeOf(model.Tags{}),
}

func ValidateUnique(db *gorm.DB, fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	tableName := getModelFromTag(fl)
	fmt.Println("for table", tableName)

	exists := UniqueExistsInTable(db, value, tableName)

	return !exists
}

func UniqueExistsInTable(db *gorm.DB, value interface{}, tableName string) bool {
	parts := strings.Split(tableName, ";")
	fmt.Printf("split: %s\n", parts)
	modelName := parts[0]
	columnName := parts[1]

	modelType, ok := modelMap[modelName]
	if !ok {
		return false
	}

	modelInstance := reflect.New(modelType).Interface()
	fmt.Printf("model instance: %s\n", modelType)

	var err error
	if len(parts) > 2 {
		idCondition := parts[2]
		err = db.Table(modelName).Where(columnName+" = ? AND "+idCondition, value).First(modelInstance).Error
	} else {
		err = db.Table(modelName).Where(columnName+" = ?", value).First(modelInstance).Error
	}
	if err != nil {
		return false
	}
	return true
}

func getModelFromTag(fl validator.FieldLevel) string {
	// Assuming 'validate' tag is in the format "unique=tableName;columnName;columnID" columnID is optional when update data
	validateTag := fl.Param()

	parts := strings.Split(validateTag, "=")
	if len(parts) >= 2 {
		return parts[1]
	}

	return validateTag
}
