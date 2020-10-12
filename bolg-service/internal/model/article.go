package model

import "github.com/jinzhu/gorm"

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_artcile"
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? and is_del = ? ", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Updates(values).Where("id = ? and is_del = ?", a.ID, a.IsDel).Error; err != nil {
		return err
	}
	return nil
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? and state = ? and is_del = ? ", a.ID, a.State, 0)
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

type ArticleRow struct {
	ArticleId     uint32
	TagId         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

func (a Article) ListByTagId(db *gorm.DB, tagid uint32, pageoffset, pagesize int) ([]*ArticleRow, error) {
	fields := []string{
		"ar.id as article_id",
		"ar.title as article_title",
		"ar.desc as article_desc",
		"ar.cover_image_url",
		"ar.content",
		"t.id as tag_id",
		"t.name as tag_name",
	}
	if pageoffset >= 0 && pagesize > 0 {
		db = db.Offset(pageoffset).Limit(pagesize)
	}
	rows, err := db.
		Select(fields).Table(ArticleTag{}.TableName()+" as at").
		Joins("left join `"+Tag{}.TableName()+"` as t on at.tag_id = t.id").
		Joins("left join `"+Article{}.TableName()+"` as ar on at.article_id = ar.id").
		Where("at.`tag_id` = ? and ar.state = ? and ar.is_del = ?", tagid, a.State, 0).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleId, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagId, &r.TagName); err != nil {
			return nil, err
		}
		articles = append(articles, r)
	}
	return articles, nil
}

func (a Article) CountByTagid(db *gorm.DB, tagid uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" as at").
		Joins("left join `"+Tag{}.TableName()+"` as t on at.tag_id =t.id ").
		Joins("left join `"+Article{}.TableName()+"` as ar on at.article_id = ar.id").
		Where(" at.`tag_id` = ? and ar.state = ? and ar.is_del = ? ", tagid, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
