package dao

import "github.com/goproject/blog-service/internal/model"

func (a *Dao) GetArticletagByAid(articleId uint32) (model.ArticleTag, error) {
	articletag := model.ArticleTag{ArticleID: articleId}
	return articletag.GetByArticleId(a.engine)
}

func (a *Dao) GetArticleListByTid(tagid uint32) ([]*model.ArticleTag, error) {
	artciletag := model.ArticleTag{TagID: tagid}
	return artciletag.ListByTid(a.engine)
}

func (a *Dao) GetArticleListByAidsi(aids []uint32) ([]*model.ArticleTag, error) {
	artciletag := model.ArticleTag{}
	return artciletag.ListByAids(a.engine, aids)
}

func (a *Dao) CreateArticleTag(articleid, tagid uint32, createby string) error {
	articletag := model.ArticleTag{
		Model: &model.Model{
			CreatedBy: createby,
		},
		ArticleID: articleid,
		TagID:     tagid,
	}
	return articletag.Create(a.engine)
}

func (a *Dao) UpdateArticleTag(articleid, tagid uint32, modifiedby string) error {
	articletag := model.ArticleTag{
		ArticleID: articleid,
	}
	values := map[string]interface{}{
		"article_id":  articleid,
		"tag_id":      tagid,
		"modified_by": modifiedby,
	}

	return articletag.UpdateOne(a.engine, values)
}

func (a *Dao) DeleteArticleTag(articleid uint32) error {
	articletag := model.ArticleTag{ArticleID: articleid}
	return articletag.DeleteOne(a.engine)
}
