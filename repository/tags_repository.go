package repository

import "gin-boilerplate/model"

type TagsRepository interface {
	Save(tags model.Tags)
	Update(tags model.Tags)
	Delete(tagsId int)
	FindById(tagsId int) (tags model.Tags, err error)
	FindAll(limit int, page int, filters map[string]interface{}) ([]model.Tags)
	Count() (int64, error)
}