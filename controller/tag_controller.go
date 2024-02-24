package controller

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/service"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TagController struct {
	tagService service.TagService
}

func NewTagController(service service.TagService) *TagController {
	return &TagController{
		tagService: service,
	}
}

// CreateTags		godoc
// @Summary			Create tags
// @Description		Save tags data in Db.
// @Param			tags body request.CreateTagsRequest true "Create tags"
// @Produce			application/json
// @Tags			tags
// @Success			200 {object} response.Response{}
// @Router			/tags [post]
// @Security        Bearer
func (controller *TagController) Create(ctx *gin.Context) {
	log.Info().Msg("create tags")
	createTagsRequest := request.CreateTagsRequest{}
	err := ctx.ShouldBindJSON(&createTagsRequest)
	helper.ErrorPanic(err)
	
	controller.tagService.Create(createTagsRequest)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Message: "Created successfully",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateTags		godoc
// @Summary			Update tags
// @Description		Update tags data.
// @Param			tagId path string true "update tags by id"
// @Param			tags body request.UpdateTagsRequest true  "Update tags"
// @Tags			tags
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/tags/{tagId} [patch]
// @Security        Bearer
func (controller *TagController) Update(ctx *gin.Context) {
	log.Info().Msg("update tags")
	updateTagsRequest := request.UpdateTagsRequest{}
	err := ctx.ShouldBindJSON(&updateTagsRequest)
	helper.ErrorPanic(err)
	
	tagId := ctx.Param("tagId")
	id, err := strconv.Atoi(tagId)
	helper.ErrorPanic(err)
	updateTagsRequest.Id = id
	
	controller.tagService.Update(updateTagsRequest)
	
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// DeleteTags		godoc
// @Summary			Delete tags
// @Param		    tagId path string true "delete tags by id"
// @Description		Remove tags data by id.
// @Produce			application/json
// @Tags			tags
// @Success			200 {object} response.Response{}
// @Router			/tags/{tagId} [delete]
// @Security        Bearer
func (controller *TagController) Delete(ctx *gin.Context) {
	log.Info().Msg("delete tags")
	tagId := ctx.Param("tagId")
	id, err := strconv.Atoi(tagId)
	helper.ErrorPanic(err)

	controller.tagService.Delete(id)
	
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindByIdTags 		godoc
// @Summary				Get Single tags by id.
// @Param				tagId path string true "update tags by id"
// @Description			Return the tahs whoes tagId value mathes id.
// @Produce				application/json
// @Tags				tags
// @Success				200 {object} response.Response{}
// @Router				/tags/{tagId} [get]
// @Security            Bearer
func (controller *TagController) FindById(ctx *gin.Context) {
	log.Info().Msg("findbyid tags")
	tagId := ctx.Param("tagId")
	id, err := strconv.Atoi(tagId)
	helper.ErrorPanic(err)
	
	tagResponse := controller.tagService.FindById(id)
	
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   tagResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAllTags 		godoc
// @Summary			Get All tags.
// @Description		Return list of tags.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    name query string false "Name"
// @Produce		    application/json
// @Tags			tags
// @Success         200 {object} response.Pagination{}
// @Router			/tags [get]
// @Security        Bearer
func (controller *TagController) FindAll(ctx *gin.Context) {
	log.Info().Msg("findAll tags")
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	page,  _ := strconv.Atoi(ctx.Query("page"))

	filters := make(map[string]interface{})
    filters["name"] = ctx.Query("name")

	if limit < 1 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}
	tagResponse, err := controller.tagService.FindAll(limit, page, filters)
	helper.ErrorPanic(err)
	
	webResponse := response.Pagination{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   tagResponse.Data,
		Meta:   tagResponse.Meta,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
	
}
