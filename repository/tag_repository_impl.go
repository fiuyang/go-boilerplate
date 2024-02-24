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

type TagRepositoryImpl struct {
	Db *gorm.DB
}

func NewTagRepositoryImpl(Db *gorm.DB) TagRepository {
	return &TagRepositoryImpl{Db: Db}
}

func (repo *TagRepositoryImpl) Delete(tagId int) {
	var tag model.Tags
	result := repo.Db.Where("id = ?", tagId).Delete(&tag)
	helper.ErrorPanic(result.Error)
}

func (repo *TagRepositoryImpl) FindById(tagsId int) (tags model.Tags, err error) {
	var tag model.Tags
	result := repo.Db.Find(&tag, tagsId)

	if result.RowsAffected == 0 {
		return tag, errors.New("tag is not found")
	}

	return tag, nil
}
	
func (repo *TagRepositoryImpl) Save(tag model.Tags) {
	result := repo.Db.Create(&tag)
	helper.ErrorPanic(result.Error)
}

func (repo *TagRepositoryImpl) Update(tag model.Tags) {
	var updateTag = request.UpdateTagsRequest{
		Id:   tag.Id,
		Name: tag.Name,
	}
	result := repo.Db.Model(&tag).Updates(updateTag)
	helper.ErrorPanic(result.Error)
}

func (repo *TagRepositoryImpl) FindAll(limit int, page int, filters map[string]interface{}) []model.Tags {
	var tags []model.Tags
	var totalRows int64

    query := repo.Db.Model(&tags)

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

func (repo *TagRepositoryImpl) Count() (int64, error) {
	var tags []model.Tags
	var totalRows int64
	result := repo.Db.Find(&tags).Count(&totalRows)
	helper.ErrorPanic(result.Error)
	return totalRows, nil
}
	