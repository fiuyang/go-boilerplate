package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"
)

type TagsServiceImpl struct {
	TagsRepository repository.TagsRepository
}

func NewTagsServiceImpl(tagRepository repository.TagsRepository) TagsService {
	return &TagsServiceImpl{
		TagsRepository: tagRepository,
	}
}


func (t *TagsServiceImpl) Create(tags request.CreateTagsRequest) {
	tagModel := model.Tags{
		Name: tags.Name,
	}
	t.TagsRepository.Save(tagModel)
}


func (t *TagsServiceImpl) Delete(tagsId int) {
	t.TagsRepository.Delete(tagsId)
}


func (t *TagsServiceImpl) FindAll(limit int, page int, filters map[string]interface{}) (response.Pagination, error) {

	paginateResponse := response.Pagination{}
	
	result := t.TagsRepository.FindAll(limit, page, filters)

	var tags []response.TagsResponse
	for _, value := range result {
		tag := response.TagsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	totalData, _ := t.TagsRepository.Count()

	paginateResponse.Data = tags
	paginateResponse.Meta.Page = page
	paginateResponse.Meta.Limit = limit
	paginateResponse.Meta.TotalData = totalData
	paginateResponse.Meta.TotalPage = totalData / int64(limit)

	return paginateResponse, nil
}



func (t *TagsServiceImpl) FindById(tagsId int) response.TagsResponse {
	tagData, err := t.TagsRepository.FindById(tagsId)
	helper.ErrorPanic(err)
	
	tagResponse := response.TagsResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}
	return tagResponse
}


func (t *TagsServiceImpl) Update(tags request.UpdateTagsRequest) {
	tagData, err := t.TagsRepository.FindById(tags.Id)
	helper.ErrorPanic(err)
	tagData.Name = tags.Name
	t.TagsRepository.Update(tagData)
}
