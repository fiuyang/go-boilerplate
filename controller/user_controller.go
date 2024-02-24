package controller

import (
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
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
	users := controller.userService.FindAll(filters)

	webResponse := response.Response{
		Code:   200,
		Status: "Ok",
		Data:   users,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			Create User
// @Description		Save user data in Db.
// @Param			users body request.CreateUsersRequest true "Create user"
// @Produce			application/json
// @Tags			users
// @Success			200 {object} response.Response{}
// @Router			/users [post]
// @Security        Bearer
func (controller *UserController) Create(ctx *gin.Context) {
	createUsersRequest := request.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)

	controller.userService.Create(createUsersRequest)

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully created user!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateTags		godoc
// @Summary			Update user
// @Description		Update user data.
// @Param			userId path string true "update user by id"
// @Param			user body request.UpdateUsersRequest true  "Update user"
// @Tags			users
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/users/{userId} [patch]
// @Security        Bearer
func (controller *UserController) Update(ctx *gin.Context) {
	log.Info().Msg("update users")
	updateUsersRequest := request.UpdateUsersRequest{}
	err := ctx.ShouldBindJSON(&updateUsersRequest)
	helper.ErrorPanic(err)

	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.ErrorPanic(err)
	updateUsersRequest.Id = id

	controller.userService.Update(updateUsersRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Successfully updated user!",
		Data:    nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// DeleteTags		godoc
// @Summary         Delete users by array of IDs
// @Description     Remove user data by providing an array of user IDs in the request body.
// @Param           request body []int true "Array of user IDs to delete"
// @Produce			application/json
// @Tags			users
// @Success			200 {object} response.Response{}
// @Router			/users/bulk/ [post]
// @Security        Bearer
func (controller *UserController) BulkDelete(ctx *gin.Context) {
	log.Info().Msg("bulk delete users")

	var userIds []int

	controller.userService.BulkDelete(userIds)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Successfully bulk delete user!",
		Data:    nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindByIdTags 		godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "update user by id"
// @Description			Return the tahs whoes userId value mathes id.
// @Produce				application/json
// @Tags				users
// @Success				200 {object} response.Response{}
// @Router				/users/{userId} [get]
// @Security            Bearer
func (controller *UserController) FindById(ctx *gin.Context) {
	log.Info().Msg("findbyid user")
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.ErrorPanic(err)

	userResponse := controller.userService.FindById(id)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   userResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// ExportTags 		    godoc
// @Summary				Export Excel User.
// @Description			Return the export excel user.
// @Produce				application/json
// @Tags				users
// @Success				200 {object} response.Response{}
// @Router				/users/export [get]
// @Security            Bearer
func (controller *UserController) Export(ctx *gin.Context) {
	filePath, err := controller.userService.Export()
	helper.ErrorPanic(err)
	defer os.Remove(filePath) // Remove the file after the function exits

	fileName := filepath.Base(filePath)
	// Set headers for the Excel file
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// Read the Excel file and write to the response body
	data, err := os.ReadFile(filePath)
	helper.ErrorPanic(err)

	// Write data to the response body
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// ImportTags 		    godoc
// @Summary				Import Excel User.
// @Description			Upload and import user data from an Excel file.
// @Produce				application/json
// @Tags				users
// @Param               file formData file true "Excel file to import"
// @Success				200 {object} response.Response{}
// @Router				/users/import [post]
// @Security            Bearer
func (controller *UserController) Import(ctx *gin.Context) {
	// Retrieve the uploaded file
	file, err := ctx.FormFile("file")
	helper.ErrorPanic(err)

	// Call the usersService method to handle the import
	err = controller.userService.Import(file)
	helper.ErrorPanic(err)

	// Return a success response
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "User data imported successfully",
		Data:    nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
