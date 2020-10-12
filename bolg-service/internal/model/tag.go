package model

import (
	"github.com/goproject/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

// tag.go
type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) TableName() string {
	return "blog_tag"
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	db = db.Where("id = ? and is_del = ? and state = ?", t.ID, 0, t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Find(&tag).Error; err != nil {
		return tag, err
	}
	return tag, nil
}

func (t Tag) List(db *gorm.DB, pageoffset, pagesize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageoffset >= 0 && pagesize > 0 {
		db = db.Offset(pageoffset).Limit(pagesize)
	}
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, value interface{}) error {
	return db.Model(&t).Where("id = ? and is_del = ?", t.ID, t.IsDel).Updates(value).Error
}

func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
