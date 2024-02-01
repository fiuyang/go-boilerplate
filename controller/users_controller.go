package controller

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/service"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	usersService service.UsersService
	Validate       *validator.Validate
}

func NewUsersController(service service.UsersService, validate *validator.Validate) *UserController {
	return &UserController{
		usersService: service,
		Validate:       validate,
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
	users := controller.usersService.FindAll(filters)

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully fetch all!",
		Data:    users,
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
	
	if err := controller.Validate.Struct(createUsersRequest); err != nil {
		if helper.ValidationError(err, ctx, createUsersRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	controller.usersService.Create(createUsersRequest)
	
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
// @Router			/user/{userId} [patch]
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
	
	controller.usersService.Update(updateUsersRequest)
	
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Message: "Successfully updated user!",
		Data:   nil,
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
	if err := ctx.ShouldBindJSON(&userIds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.usersService.BulkDelete(userIds)
	
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Message: "Successfully bulk delete user!",
		Data:   nil,
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
	
	userResponse := controller.usersService.FindById(id)
	
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
	filePath, err := controller.usersService.Export()
	helper.ErrorPanic(err)
    // Set headers for the Excel file
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=user.xlsx")

	// Read the Excel file and write to the response body
	data, err := ioutil.ReadFile(filePath)
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
	err = controller.usersService.Import(file)
	helper.ErrorPanic(err)

	// Return a success response
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Message: "User data imported successfully",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}