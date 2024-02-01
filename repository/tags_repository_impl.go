package repository

import (
	"errors"
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"

	"gorm.io/gorm"
)

type TagsRepositoryImpl struct {
	Db *gorm.DB
}

func NewTagsREpositoryImpl(Db *gorm.DB) TagsRepository {
	return &TagsRepositoryImpl{Db: Db}
}

func (t *TagsRepositoryImpl) Delete(tagsId int) {
	var tags model.Tags
	result := t.Db.Where("id = ?", tagsId).Delete(&tags)
	helper.ErrorPanic(result.Error)
}

func (t *TagsRepositoryImpl) FindById(tagsId int) (tags model.Tags, err error) {
	var tag model.Tags
	result := t.Db.Find(&tag, tagsId)
	if result != nil {
		return tag, nil
	} else {
		return tag, errors.New("tag is not found")
	}
}
	
func (t *TagsRepositoryImpl) Save(tags model.Tags) {
	result := t.Db.Create(&tags)
	helper.ErrorPanic(result.Error)
}

func (t *TagsRepositoryImpl) Update(tags model.Tags) {
	var updateTag = request.UpdateTagsRequest{
		Id:   tags.Id,
		Name: tags.Name,
	}
	result := t.Db.Model(&tags).Updates(updateTag)
	helper.ErrorPanic(result.Error)
}

func (t *TagsRepositoryImpl) FindAll(limit int, page int, filters map[string]interface{}) []model.Tags {
	var tags []model.Tags
	var totalRows int64

    query := t.Db.Model(&tags)

    // Menambahkan filter ke dalam query
    for field, value := range filters {
        // Memeriksa tipe data dan apakah nilainya tidak kosong atau nol
        switch v := value.(type) {
        case string:
            if v != "" {
                query = query.Where(fmt.Sprintf("%s = ?", field), v)
            }
        case int:
            if v != 0 {
                query = query.Where(fmt.Sprintf("%s = ?", field), v)
            }
        }
    }

    result := query.Scopes(response.Scopes(page, limit)).Find(&tags).Count(&totalRows)
    helper.ErrorPanic(result.Error)

    return tags
}

func (t *TagsRepositoryImpl) Count() (int64, error) {
	var tags []model.Tags
	var totalRows int64
	result := t.Db.Find(&tags).Count(&totalRows)
	helper.ErrorPanic(result.Error)
	return totalRows, nil
}
	