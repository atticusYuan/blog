package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int

	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	return tag.ID > 0, err
}
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? ", id).First(&tag).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return tag.ID > 0, nil
}

func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}

func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}
func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error
	if err != nil {
		return err
	}
	return err
}

func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}
