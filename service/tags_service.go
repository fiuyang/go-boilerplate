package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
)

type TagsService interface {
	Create(tags request.CreateTagsRequest)
	Update(tags request.UpdateTagsRequest)
	Delete(tagsId int)
	FindById(tagsId int) response.TagsResponse
	FindAll(limit int, page int, filters map[string]interface{}) (response.Pagination, error)
}
