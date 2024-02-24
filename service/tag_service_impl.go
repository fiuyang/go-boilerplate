package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/exception"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"

	"github.com/go-playground/validator/v10"
)

type TagServiceImpl struct {
	TagRepository repository.TagRepository
	Validate      *validator.Validate
}

func NewTagServiceImpl(tagRepository repository.TagRepository, validate *validator.Validate) TagService {
	return &TagServiceImpl{
		TagRepository: tagRepository,
		Validate:      validate,
	}
}


func (service *TagServiceImpl) Create(tag request.CreateTagsRequest) {
	err := service.Validate.Struct(tag)
	helper.ErrorPanic(err)

	dataset := model.Tags{
		Name: tag.Name,
	}
	service.TagRepository.Save(dataset)
}


func (service *TagServiceImpl) Delete(tagId int) {
	service.TagRepository.Delete(tagId)
}


func (service *TagServiceImpl) FindAll(limit int, page int, filters map[string]interface{}) (response.Pagination, error) {

	paginateResponse := response.Pagination{}
	
	result := service.TagRepository.FindAll(limit, page, filters)

	var tags []response.TagsResponse
	for _, value := range result {
		tag := response.TagsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	totalData, _ := service.TagRepository.Count()

	paginateResponse.Data = tags
	paginateResponse.Meta.Page = page
	paginateResponse.Meta.Limit = limit
	paginateResponse.Meta.TotalData = totalData
	paginateResponse.Meta.TotalPage = totalData / int64(limit)

	return paginateResponse, nil
}



func (service *TagServiceImpl) FindById(tagId int) response.TagsResponse {
	dataset, err := service.TagRepository.FindById(tagId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	
	Response := response.TagsResponse{
		Id:   dataset.Id,
		Name: dataset.Name,
	}
	return Response
}


func (service *TagServiceImpl) Update(tag request.UpdateTagsRequest) {
	err := service.Validate.Struct(tag)
	helper.ErrorPanic(err)

	dataset, err := service.TagRepository.FindById(tag.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	dataset.Name = tag.Name
	service.TagRepository.Update(dataset)
}
