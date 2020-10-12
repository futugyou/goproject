package service

import (
	"github.com/goproject/blog-service/internal/dao"
	"github.com/goproject/blog-service/internal/model"
	"github.com/goproject/blog-service/pkg/app"
)

type Article struct {
	Id            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
}
type ArticleRequest struct {
	Id    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

func (svc *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := svc.dao.GetArticle(param.Id, param.State)
	if err != nil {
		return nil, err
	}
	articletag, err := svc.dao.GetArticletagByAid(article.ID)
	if err != nil {
		return nil, err
	}
	tag, err := svc.dao.GetTag(articletag.TagID, param.State)

	if err != nil {
		return nil, err
	}
	return &Article{
		Id:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tag:           &tag,
	}, nil
}

type ArticleListRequest struct {
	TagId uint32 `form:"tag_id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

func (svc *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	articleCount, err := svc.dao.CountArtileListbyTagid(param.TagId, param.State)
	if err != nil {
		return nil, 0, err
	}
	articles, err := svc.dao.GetArticleListByTagid(param.TagId, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}
	var articleList []*Article
	for _, article := range articles {
		articleList = append(articleList, &Article{
			Id:            article.ArticleId,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tag: &model.Tag{
				Model: &model.Model{ID: article.TagId},
				Name:  article.TagName},
		})
	}
	return articleList, articleCount, nil
}

type CreateArticleRequest struct {
	TagId uint32 `form:"tag_id" binding:"required,gte=1"`
	//Name          string `form:"name" binding:"required,min=3,max=100"`
	CreatedBy     string `form:"created_by" binding:"required,min=3,max=100"`
	Title         string `form:"title" binding:"required,min=3,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
	Desc          string `form:"desc" binding:"min=3,max=100"`
	Content       string `form:"content" binding:"min=3,max=100"`
	CoverImageUrl string `form:"cover_image_url" binding:"min=3,max=100"`
}

func (svc *Service) CreateArticle(param *CreateArticleRequest) error {
	article, err := svc.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		CreatedBy:     param.CreatedBy,
	})
	if err != nil {
		return err
	}
	err = svc.dao.CreateArticleTag(article.ID, param.TagId, param.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

type UpdateArticleRequest struct {
	TagId         uint32 `form:"tag_id" binding:"required,gte=1"`
	ID            uint32 `form:"id" binding:"required,gte=1"`
	Name          string `form:"name" binding:"required,min=3,max=100"`
	ModifiedBy    string `form:"modified_dy" binding:"required,min=3,max=100"`
	Title         string `form:"title" binding:"required,min=3,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
	Desc          string `form:"desc" binding:"min=3,max=100"`
	Content       string `form:"conten" binding:"min=3,max=100"`
	CoverImageUrl string `form:"cover_image_url" binding:"min=3,max=100"`
}

func (svc *Service) UpdateArticle(param *UpdateArticleRequest) error {
	err := svc.dao.UpdateArticle(&dao.Article{
		Id:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	})
	if err != nil {
		return err
	}
	err = svc.dao.UpdateArticleTag(param.ID, param.TagId, param.ModifiedBy)
	if err != nil {
		return err
	}
	return nil
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

func (svc *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := svc.dao.DeleteTag(param.ID)
	if err != nil {
		return err
	}
	err = svc.dao.DeleteArticleTag(param.ID)
	if err != nil {
		return err
	}
	return nil
}
