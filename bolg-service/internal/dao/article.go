package dao

import (
	"github.com/goproject/blog-service/internal/model"
	"github.com/goproject/blog-service/pkg/app"
)

type Article struct {
	Id            uint32 `json:"id"`
	TagId         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

func (a *Dao) CreateArticle(param *Article) (*model.Article, error) {
	article := model.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreatedBy},
	}
	return article.Create(a.engine)
}
func (a *Dao) UpdateArticle(param *Article) error {
	article := model.Article{Model: &model.Model{ID: param.Id}}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.CoverImageUrl != "" {
		values["cover_image_url"] = param.CoverImageUrl
	}
	if param.Desc != "" {
		values["dontent"] = param.Desc
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	return article.Update(a.engine, values)
}

func (a *Dao) GetArticle(id uint32, state uint8) (model.Article, error) {
	article := model.Article{Model: &model.Model{ID: id}, State: state}
	return article.Get(a.engine)
}

func (a *Dao) DeleteArticle(id uint32) error {
	article := model.Article{Model: &model.Model{ID: id}}
	return article.Delete(a.engine)
}

func (a *Dao) CountArtileListbyTagid(id uint32, state uint8) (int, error) {
	article := model.Article{State: state}
	return article.CountByTagid(a.engine, id)
}

func (a *Dao) GetArticleListByTagid(id uint32, state uint8, page, pagesize int) ([]*model.ArticleRow, error) {
	article := model.Article{State: state}
	return article.ListByTagId(a.engine, id, app.GetPageOffset(page, pagesize), pagesize)
}
