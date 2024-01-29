package controller

import (
	"gin-boilerplate/data/response"
	"gin-boilerplate/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usersService service.UsersService
}

func NewUsersController(service service.UsersService) *UserController {
	return &UserController{usersService: service}
}

// FindAllTags 		godoc
// @Summary			Get All users.
// @Description		Return list of users.
// @Param		    email query string false "Email"
// @Param		    username query string false "Username"
// @Produce		    application/json
// @Tags			users
// @Success         200 {object} response.Response{}
// @Router			/users [get]
// @Security        Bearer
func (controller *UserController) GetUsers(ctx *gin.Context) {
	filters := make(map[string]string)
    filters["username"] = ctx.Query("username")
    filters["email"] = ctx.Query("email")
	users := controller.usersService.FindAll(filters)

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully fetch all!",
		Data:    users,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
