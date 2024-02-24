package exception

import (
	"fmt"
	"gin-boilerplate/data/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ErrorHandlers(ctx *gin.Context, err interface{}) {
	if validationError(ctx, err) {
		return
	}

	if notFoundError(ctx, err) {
		return
	}

	internalServerError(ctx, err)
}

func validationError(ctx *gin.Context, err interface{}) bool {
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		report := make(map[string]string)

		for _, e := range castedObject {
			fieldName := e.Field()
			switch e.Tag() {
			case "required":
				report[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				report[fieldName] = fmt.Sprintf("%s is not valid email", fieldName)
			case "gte":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "lte":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "unique":
				report[fieldName] = fmt.Sprintf("%s has already been taken %s", fieldName, e.Param())
			case "max":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "min":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "numeric":
				report[fieldName] = fmt.Sprintf("%s value must be numeric", fieldName)
			}
		}

		ctx.JSON(http.StatusBadRequest, response.Error{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Errors: report,
		})
		return true
	}
	return false
}

func notFoundError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		ctx.JSON(http.StatusNotFound, response.Error{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Errors: exception.Error,
		})
		return true
	}
	return false
}

// func internalServerError(ctx *gin.Context, err interface{}) {
// 	ctx.JSON(http.StatusInternalServerError, response.Error{
// 		Code:   http.StatusInternalServerError,
// 		Status: "INTERNAL SERVER ERROR",
// 		Errors: err,
// 	})
// }

func internalServerError(ctx *gin.Context, err interface{}) {
    var errMsg string
    if err != nil {
        errMsg = fmt.Sprintf("%v", err)
    } else {
        errMsg = "Unknown error occurred"
    }

    ctx.JSON(http.StatusInternalServerError, response.Error{
        Code:   http.StatusInternalServerError,
        Status: "INTERNAL SERVER ERROR",
        Errors: errMsg,
    })
}