package model

import (
	"github.com/jinzhu/gorm"
)

type ArticleTag struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (a ArticleTag) GetByArticleId(db *gorm.DB) (ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? and is_del = ? ", a.ArticleID, 0).
		First(&articleTag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return articleTag, err
	}
	return articleTag, nil
}

func (a ArticleTag) ListByTid(db *gorm.DB) ([]*ArticleTag, error) {
	var articleTag []*ArticleTag
	if err := db.Where("tag_id = ? and is_del = ?", a.TagID, 0).Find(&articleTag).Error; err != nil {
		return nil, err
	}
	return articleTag, nil
}

func (a ArticleTag) ListByAids(db *gorm.DB, aids []uint32) ([]*ArticleTag, error) {
	var articleTag []*ArticleTag
	err := db.Where("article_id in (?) and is_del = ?", aids, 0).Find(&articleTag).Error
	if err != nil {
		return nil, err
	}
	return articleTag, nil
}

func (a ArticleTag) Create(db *gorm.DB) error {
	if err := db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) UpdateOne(db *gorm.DB, value interface{}) error {
	err := db.Model(&a).Where("article_id = ? and is_del = ?", a.ArticleID, 0).Limit(1).Updates(value).Error
	if err != nil {
		return err
	}
	return nil
}
func (a ArticleTag) Delete(db *gorm.DB) error {
	err := db.Where("id = ? and is_del = ?", a.ID, 0).Delete(&a).Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) DeleteOne(db *gorm.DB) error {
	err := db.Where("article_id = ? and is_del = ?", a.ArticleID, 0).Delete(&a).Limit(1).Error
	if err != nil {
		return err
	}
	return nil
}
