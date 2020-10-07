package model

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	*Model
	TageID    uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (a ArticleTag) GetByArticleId(db *gorm.DB, id uint32) (ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? and is_del = ? ", a.TageID, 0).
		First(&articleTag).Error
	if err != nil {
		return articleTag, err
	}
	return articleTag, nil
}
