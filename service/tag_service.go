package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
)

type TagService interface {
	Create(tag request.CreateTagsRequest)
	Update(tag request.UpdateTagsRequest)
	Delete(tagId int)
	FindById(tagId int) response.TagsResponse
	FindAll(limit int, page int, filters map[string]interface{}) (response.Pagination, error)
}
