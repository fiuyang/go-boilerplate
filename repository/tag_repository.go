package repository

import "gin-boilerplate/model"

type TagRepository interface {
	Save(tag model.Tags)
	Update(tag model.Tags)
	Delete(tagId int)
	FindById(tagId int) (tags model.Tags, err error)
	FindAll(limit int, page int, filters map[string]interface{}) ([]model.Tags)
	Count() (int64, error)
}